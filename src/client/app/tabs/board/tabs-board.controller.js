angular.module("scoreboard").controller("BoardCtrl", function BoardCtrl (Parameters, toaster, $q) {
    "use strict";
    var self = this;

    this.getCurrentParameters = function () {
        Parameters.read().$promise.then(function (parameters) {
            self.parameters = parameters;
            self.parameters.game.teamA.textColor = getLum(parameters.game.teamA.color_code) > .5 ? "#000000" : "#ffffff";
            self.parameters.game.teamB.textColor = getLum(parameters.game.teamB.color_code) > .5 ? "#000000" : "#ffffff";
            self.parameters.game.jamScoreA = 0;
            self.parameters.game.jamScoreB = 0;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Game", body:"Cannot start game. No parameter is set"});
            self.errors.push({
                missingParameters: true
            })
            return $q.reject(err);
        });
    };

    function getLum(str){
        var matcher = str.match(/^#([0-9A-Fa-f]{2})([0-9A-Fa-f]{2})([0-9A-Fa-f]{2})$/);
        if (matcher) {
            var r = parseInt(matcher[1], 16) / 255
            var g = parseInt(matcher[2], 16) / 255;
            var b = parseInt(matcher[3], 16) / 255;
            var max = Math.max(r, g, b), min = Math.min(r, g, b);
            return (max + min) / 2;
        }
    }

    this.saveScore = function () {
        console.log("saving score", this.parameters.game.scoreA, this.parameters.game.scoreB);
    };

    this.adjustScore = function (a, b) {
        self.parameters.game.jamScoreA += a;
        self.parameters.game.jamScoreB += b;
        self.parameters.game.scoreA += a;
        self.parameters.game.scoreB += b;
        this.saveScore();
    }

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.errors = [];
        this.getCurrentParameters();
    };

});