angular.module("scoreboard").controller("TeamDetailCtrl", function TeamDetailCtrl ($q, $scope, $state, $stateParams, Team, Player, FileUploader, toaster) {
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
            self.teamTitle = team.name;
            return team;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Team", body:"could not be found"});
            return $q.reject(err);
        });
    };

    /**
     * Get all the players
     * @return {Promise}
     */
    this.getPlayers = function () {
        return Player.list({
            teamId: $stateParams.teamId
        }).$promise.then(function (players) {
            self.players = players;
            return players;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Players", body:"could not be found"});
            return $q.reject(err);
        });
    };

    /**
     * Delete the current team
     * @return {Promise}
     */
    this.deleteTeam = function () {
        this.teamBusy = true;
        return Team.delete({
            teamId: $stateParams.teamId
        }).$promise.then(function () {
            $scope.$emit("refresh-team-list");
            toaster.pop({ type: "success", title: "Team deleted"});
            return $state.go("main.teams");
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Team", body:"could not be deleted"});
            return $q.reject(err);
        }).finally(function() {
            delete self.teamBusy;
        });
    }

    /**
     * Save the current team
     * @return {Promise}
     */
    this.saveTeam = function () {
        this.teamBusy = true;
        return Team.update({
            teamId: $stateParams.teamId
        }, {
            name: this.team.name,
            color: this.team.color,
            color_code: this.team.color_code
        }).$promise.then(function (team) {
            $scope.$emit("refresh-team-list");
            self.teamTitle = team.name;
            toaster.pop({ type: "success", title: "Team saved"});
            return team;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Team", body:"could not be saved"});
            return $q.reject(err);
        }).finally(function() {
            delete self.teamBusy;
        });
    }

    /**
     * Add a player
     * @return {Promise}
     */
    this.addPlayer = function () {
        this.playerAdding = true;
        return Player.create({
            teamId: $stateParams.teamId
        }, {
            name: this.newPlayer.name,
            number: this.newPlayer.number
        }).$promise.then(function (player) {
            self.newPlayer = null;
            self.players.push(player);
            self.players = _.sortBy(self.players, ["number"]);
            return player;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Player", body:"could not be created"});
            return $q.reject(err);
        }).finally(function() {
            delete self.playerAdding;
        });
    }

    /**
     * Delete a player
     * @param {Object} player Player to save 
     * @return {Promise}
     */
    this.deletePlayer = function (player) {
        player.deleting = true;
        return Player.delete({
            teamId: $stateParams.teamId,
            playerId: player.id
        }).$promise.then(function (result) {
            _.remove(self.players, { id: player.id });
            return result;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Player", body:"could not be deleted"});
            return $q.reject(err);
        }).finally(function() {
            delete player.deleting;
        });
    };

    /**
     * Save a player
     * @param {Object} player Player to save 
     * @return {Promise}
     */
    this.savePlayer = function (player) {
        player.saving = true;
        return Player.update({
            teamId: $stateParams.teamId,
            playerId: player.id
        }, {
            name: player.name,
            number: player.number
        }).$promise.then(function (result) {
            self.players = _.sortBy(self.players, ["number"]);
            toaster.pop({ type: "success", title: "Player saved"});
            return result;
        }).catch(function (err) {
            toaster.pop({ type: "error", title: "Player", body:"could not be saved"});
            return $q.reject(err);
        }).finally(function() {
            delete player.saving;
        });
    };

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.uploader = new FileUploader({
            url: "/teams/" + $stateParams.teamId + "/logo",
            queueLimit: 1,
            alias: "logo",
            autoUpload: true,
            onSuccessItem : function(fileItem, response, status, headers) {
                self.team.logo = response.logo;
                this.queue.forEach(function (item) {
                    item.remove();
                });
                toaster.pop({ type: "success", title: "Logo saved"});
            },
            onErrorItem: function(fileItem, response, status, headers) {
                this.queue.forEach(function (item) {
                    item.remove();
                });
                toaster.pop({ type: "error", title: "Logo", body:"could not be saved"});
            }
        });
        return $q.all([
            this.getTeamDetail(),
            this.getPlayers()
        ]);
    };
});