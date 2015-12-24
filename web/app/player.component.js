(function(app) {
  console.log(app)
  app.PlayerComponent = ng.core
    .Component({
      selector: 'player',
      templateUrl: 'views/player.html' 
    })
    .Class({
      constructor: function() {}
    });
})(window.app || (window.app = {}));
