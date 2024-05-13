// 独自拡張

document.addEventListener("webviewerloaded", () => {
    PDFViewerApplication.initializedPromise.then(() => {
        PDFViewerApplication.eventBus.on("pagechanging", (e) => {
            console.log('pagechanging, from ' + e.previous + ' to ' + e.pageNumber);
        });
    });
});

function addBookmark() {
    const zoom = 150;
    const page = PDFViewerApplication.pdfViewer.currentPageNumber;
    var hash = window.location.hash;
    window.location.hash = `zoom=${zoom}&page=${page}`;
};
