(function(app) {
  document.addEventListener('DOMContentLoaded', function() {
    ng.platform.browser.bootstrap(app.PlayListComponent);
 
    app.ws = new WebSocket("ws://192.168.43.28:9000/pl/pl/subscribe");
    app.send = function(message) {
      console.log("Sending: " + message);
    }

    app.ws.onmessage = function(e) {
      console.log(e.data);
    }
    

  });
})(window.app || (window.app = {}));
