(function(app) {
  app.PlayListComponent = ng.core
    .Component({
      selector: 'playlist',
      directives: [app.PlayerComponent],
      templateUrl: 'views/playlist.html'
    })

    .Class({
      constructor: function() {}
    });
})(window.app || (window.app = {}));
