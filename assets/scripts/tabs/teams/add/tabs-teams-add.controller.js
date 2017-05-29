angular.module("scoreboard").controller("TeamAddCtrl", function TeamDetailCtrl ($q, $scope, $state, $stateParams, Team, toaster) {
    "use strict";

    var self = this;

    /**
     * Add a new Team
     * @return {Promise}
     */
    this.addTeam = function () {
        return Team.create({
            name: this.team.name,
            color: this.team.color,
            color_code: this.team.color_code
        }).$promise.then(function (team) {
            $scope.$emit("refresh-team-list");
            toaster.pop({ type: "success", title: "New team"});
            return $state.go("main.teams.detail", { teamId: team.id });
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Team", body:"could not be created"});
            return $q.reject(err);
        });
    };

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.team = {
            colorCode: "#263238"
        }
    };
});