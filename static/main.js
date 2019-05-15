$(document).ready(function() {
    var x = $("div.new")
    x.each(function(i,z) {
        var group = $(z).attr('group')
        $(z).dropzone({
            url: "/post/"+group,
            forceChunking: true,
            createImageThumbnails: false,
            acceptedFiles: "image/*"
        })
    })
})