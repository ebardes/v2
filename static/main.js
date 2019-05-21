$(document).ready(function() {
    Dropzone.options.filedrop = {
        init: function () {
            this.on("queuecomplete", function (file) {
                alert("All files have uploaded ");
            });
        }
    }

    var x = $("div.new")
    x.each(function(i,z) {
        var group = $(z).attr('group')
        $(z).dropzone({
            url: "/post/"+group,
            forceChunking: true,
            createImageThumbnails: false,
            acceptedFiles: "image/*, video/*",
            queuecomplete: function(file) {
                console.log("All Done")
                location.reload();
            },
        })
    })

    x = $("div.delete")
    x.on("click", function(e) {
        $.ajax({
            url: e.target.id,
            success: function(x) {
                var p = $(e.target).parent()
                p = p.parent()
                p.remove()
            },
            error: function(e) {
                console.log(e)
            }
        })
    })
})