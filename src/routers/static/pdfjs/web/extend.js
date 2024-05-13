// 独自拡張

document.addEventListener("webviewerloaded", () => {
    PDFViewerApplication.initializedPromise.then(() => {
        PDFViewerApplication.eventBus.on("pagechanging", (e) => {
            // const msg = 'pagechanging, from ' + e.previous + ' to ' + e.pageNumber;
            updatePageGauge();
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
    var hash = window.location.hash;
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
