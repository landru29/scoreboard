angular.module("scoreboard").controller("TeamAddCtrl", function TeamDetailCtrl ($scope, $state, $stateParams, Team) {
    "use strict";

    var self = this;

    /**
     * Add a new Team
     * @return {Promise}
     */
    this.addTeam = function () {
        return Team.create({
            name: this.team.name,
            color: this.team.color
        }).$promise.then(function (team) {
            $scope.$emit("refresh-team-list");
            return $state.go("main.teams.detail", { teamId: team.id });
        });
    };


    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
    };
});