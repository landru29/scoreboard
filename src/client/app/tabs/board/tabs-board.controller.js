angular.module("scoreboard").controller("BoardCtrl", function BoardCtrl (Parameters, toaster, $q) {
    "use strict";
    var self = this;

    this.getCurrentParameters = function () {
        Parameters.read().$promise.then(function (parameters) {
            self.parameters = parameters;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Game", body:"Cannot start game. No parameter is set"});
            self.errors.push({
                missingParameters: true
            })
            return $q.reject(err);
        });
    };

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.errors = [];
        this.getCurrentParameters();
    };

});