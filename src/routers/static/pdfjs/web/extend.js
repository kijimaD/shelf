// 独自拡張

document.addEventListener("webviewerloaded", () => {
    PDFViewerApplication.initializedPromise.then(() => {
        PDFViewerApplication.eventBus.on("pagechanging", (e) => {
            console.log('pagechanging, from ' + e.previous + ' to ' + e.pageNumber);
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
