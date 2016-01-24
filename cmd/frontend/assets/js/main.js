// One file is bad, need to refactor into require.js

var UpdateTimeout;
var Token;
var Template;

function onLoad() {
    gapi.auth2.getAuthInstance().then(function(){
        onLoginChange(gapi.auth2.getAuthInstance().isSignedIn.get())
        gapi.auth2.getAuthInstance().isSignedIn.listen(onLoginChange)
    })

    var source = $("#template__subscription").html().replace(/\[\[/g, "{{").replace(/\]\]/g, "}}")
    Template = Handlebars.compile(source);

    $("#subscriptions").on("click", "a.markasread", onMarkAsRead)
    $("#subscriptions").on("click", "a.markallasread", onMarkAllAsRead)
    $(".header__float__subscribe__form").submit(onSubscribeSubmit)
    $(".subscription__confirm__body__form").submit(onSubscribeConfimSubmit)
}

function onSubscribeSubmit(evt) {
    evt.preventDefault()
    var val = $(".header__float__subscribe__form__input").val()
    Subscribe(val)
}

function onSubscribeConfimSubmit(evt) {
    evt.preventDefault()
    $(".subscription__confirm").hide()
    var val = $("#subscription__confirm__url").val()
    var title = $("#subscription__confirm__title").val()
    Subscribe(val, title)
}

function onLoginChange(loggedIn) {
    if(loggedIn){
        var googleUser = gapi.auth2.getAuthInstance().currentUser.get()
        Token = googleUser.getAuthResponse().id_token
        $(".user--valid").show()
        $(".user--anonymous").hide()
        Update()
    }else{
        $(".user--valid").hide()
        $(".user--anonymous").show()
        clearTimeout(UpdateTimeout)
    }
}

function onMarkAsRead(evt) {
    evt.preventDefault()
    var article = $(this).data("article")
    Read(article)
    $(this).closest(".subscription__info__article").hide()
}

function onMarkAllAsRead(evt) {
    evt.preventDefault()
    var blog = $(this).data("blog")
    ReadAll(blog)
    $(this).closest(".subscription__info").find(".subscription__info__article").hide()
}

function Read(article) {
    var url = window.API + "/read"
    var data = {token:Token, url: article}
    $.post(url, data, onRead)
}

function onRead(data){
    Update()
}

function ReadAll(blog) {
    var url = window.API + "/readall"
    var data = {token:Token, url: blog}
    $.post(url, data, onReadAll)
}

function onReadAll(data){
    Update()
}

function Update() {
    var url = window.API + "/updates"
    var data = {token:Token}
    $.post(url, data, onUpdate)
}

function onUpdate(data) {

    var context = {
        count: 0,
        subscriptions: []
    }

    function ninePlus(n) {
        if(n < 10) {
            return n
        }
        return "9+"
    }

    if (data.subscriptions != null) {
        if(data.updates == null){
            data.updates = []
        }
        var context = {
            count: ninePlus(data.updates.length),
            subscriptions: []
        }
        for(var i = 0; i < data.subscriptions.length; i++){
            context.subscriptions[i] = {
                url: data.subscriptions[i].url,
                title: data.subscriptions[i].title,
                count: ninePlus(data.subscriptions[i].updates.length),
                updates: data.subscriptions[i].updates
            }
        }
    }

    var html = Template(context)
    $("#subscriptions").html(html)

    $(".subscription__info__body").mCustomScrollbar({
        theme: "rounded-dark",
    });

    $("a.lb").fancybox({
        hideOnContentClick: false,
        type: "iframe",
        width: "90%",
        height: "90%",
    });

    clearTimeout(UpdateTimeout)
    UpdateTimeout = setTimeout(Update, 300000)
}

function Subscribe(blog, title) {
    if(title == undefined){
        title = ""
    }

    var url = window.API + "/subscribe"
    var data = {token:Token, url: blog, title: title}
    $.post(url, data, onSubscribe)
}

function onSubscribe(data) {
    if(!data.success && !data.error) {
        $(".subscription__confirm").show()
        $("#subscription__confirm__url").val(data.url)
        $("#subscription__confirm__title").val(data.title)
    }else{
        $(".header__float__subscribe__form__input").val("")
        Update()
    }
}

function Unsubscribe(blog) {
    var url = window.API + "/unsubscribe"
    var data = {token:Token, url: blog}
    $.post(url, data, onSubscribe)
}

function onUnsubscribe() {
    Update()
}
