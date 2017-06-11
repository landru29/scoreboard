angular.module("scoreboard").controller("GameCtrl", function GameCtrl ($scope, Game) {
    "use strict";

    var self = this;

    /**
     * Load all teams
     * @return {Promise}
     */
    this.loadGames = function () {
        return Game.list().$promise.then(function (games) {
            self.games = games;
            return games;
        });
    }

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        $scope.$on("refresh-game-list", function () {
            self.loadGames();
        });
        return this.loadGames();
    };

});