angular.module("scoreboard").controller("TeamDetailCtrl", function TeamDetailCtrl ($scope, $state, $stateParams, Team) {
    "use strict";

    var self = this;

    /**
     * Get the detail of a team
     * @return {Promise}
     */
    this.getTeamDetail = function () {
        return Team.detail({
            teamId: $stateParams.teamId
        }).$promise.then(function (team) {
            self.team = team;
            return team;
        });
    };

    /**
     * Delete the curren team
     * @return {Promise}
     */
    this.deleteTeam = function () {
        return Team.delete({
            teamId: $stateParams.teamId
        }).$promise.then(function () {
            $scope.$emit("refresh-team-list");
            return $state.go("main.teams");
        });
    }

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        return this.getTeamDetail();
    };
});