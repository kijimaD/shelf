// 独自拡張

let lastTime = Date.now();

document.addEventListener("webviewerloaded", () => {
    PDFViewerApplication.initializedPromise.then(() => {
        PDFViewerApplication.eventBus.on("pagechanging", (e) => {
            updatePageGauge();

            const lapseTime = Math.round((Date.now() - lastTime) / 1000);
            const lapseTimePad = String(lapseTime).padStart(4, '0');
            const today = new Date();
            lastTime = Date.now();

            const msg =  today.toLocaleTimeString("ja-JP") + ' [' + lapseTimePad + 's] ' + e.previous + 'p -> ' + e.pageNumber + 'p';
            let cssName = "one";
            if ((PDFViewerApplication.pdfViewer.currentPageNumber % 5) == 0) {
                cssName = "five"
            }
            if ((PDFViewerApplication.pdfViewer.currentPageNumber % 10) == 0) {
                cssName = "ten"
            }
            if ((PDFViewerApplication.pdfViewer.currentPageNumber % 100) == 0) {
                cssName = "hundred"
            }
            const content = { message: msg, cssName: cssName };
            setupToast(content);
        });
    });
});

// URLフラグメントを更新する
function addBookmark() {
    const parsedHash = new URLSearchParams(
        window.location.hash.substring(1)
    );
    const zoom = parsedHash.get("zoom");
    const page = PDFViewerApplication.pdfViewer.currentPageNumber;
    window.location.hash = `zoom=${zoom}&page=${page}`;
};

function updatePageGauge() {
    const pageGauge = document.getElementById('pageGauge');
    const cur = PDFViewerApplication.pdfViewer.currentPageNumber;
    const all = PDFViewerApplication.pdfDocument.numPages;
    const ratio = cur / all;
    const right = '#'.repeat(20 * ratio);
    const left = '_'.repeat(20 * (1 - ratio))
    const content = right + '@' + left;
    pageGauge.textContent = content;
}

// 参考: https://ics.media/entry/230530/
const setupToast = ({ message, cssName }) => {
    // トーストをDOMに追加する
    const toast = createToastElm(message, cssName);
    document.body.appendChild(toast);
    // showPopoverメソッドで表示する
    toast.showPopover();

    // setTimeoutで一定時間経ったら自動的にポップオーバーを消す
    const timer = setTimeout(() => removeToast(toast), 1000 * 60 * 10);
    // timeoutを解除するためのtimerをdataset要素として設定する
    toast.dataset.timer = timer;

    // トーストの表示時と非表示時に並び替える
    toast.addEventListener("toggle", (event) => {
        alignToast(event.newState === "closed");
    });
};

/**
 * トーストを作成します。
 * @param {string} message 表示するメッセージ
 * @param {string} cssName cssのクラス名
 * @return {HTMLDivElement} 作成したトーストエレメント
 */
const createToastElm = (message, cssName) => {
    const toast = document.createElement("div");
    toast.popover = "manual";
    toast.classList.add("toast", cssName);
    // コンテンツ
    const content = document.createElement("p");
    content.textContent = message;
    content.classList.add("toast-content");
    toast.appendChild(content);
    return toast;
};

const alignToast = (withMoveAnim) => {
    const toasts = document.querySelectorAll(".toast");
    // トーストを順番に縦に並べる
    // withMoveAnimがtrue：opacityとtranslateのアニメーション
    // withMoveAnimがfalse：opacityのアニメーション
    toasts.forEach((toast, index) => {
        toast.style.transition = withMoveAnim
                               ? "translate 0.2s linear, opacity 0.2s linear"
                               : "opacity 0.2s linear";
        toast.style.translate = `0px ${1 * index}em`;
        toast.style.opacity = 1;
    });
};

/**
 * トーストを削除します。
 * @param {HTMLDivElement} toast 削除したいトースト
 */
const removeToast = (toast) => {
    // hidePopoverメソッドで非表示にする
    toast.hidePopover();
    // 非表示にした後にDOMから削除する
    toast.remove();
    // setTimeoutを解除する
    clearTimeout(toast.dataset.timer);
};
