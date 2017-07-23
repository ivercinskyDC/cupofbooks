$(document).ready(function() {
    $.get("/user").then(function(data) {
        debugger;
        $("#username").text(data.data.username)
    }, function(err) {
        alert(err)
    })
})