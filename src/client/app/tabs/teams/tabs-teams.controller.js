angular.module("scoreboard").controller("TeamsCtrl", function TeamsCtrl ($scope, Team) {
    "use strict";

    var self = this;

    /**
     * Load all teams
     * @return {Promise}
     */
    this.loadTeams = function () {
        return Team.list().$promise.then(function (teams) {
            self.teams = teams;
            return teams;
        });
    }

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        $scope.$on("refresh-team-list", function () {
            self.loadTeams();
        });
        return this.loadTeams();
    };

});