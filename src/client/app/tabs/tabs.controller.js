angular.module("scoreboard").controller("TabsCtrl", function TabsCtrl ($scope, $transitions, TABS, $state) {
    "use strict";
    var self = this;

    this.go = function (state) {
        return $state.go(state);
    }

    $transitions.onSuccess({}, function () {
        self.setTab($state.current);
    });

    this.setTab = function (state) {
        var tabName;
        var matcher = state.name.match(/^(main\.tabs\.[\w]*)/);
        if (matcher) {
            tabName = matcher[1];
        }
        var tabIndex = _.findIndex(TABS.tabs, { state: tabName });
        if (tabIndex > -1) {
            this.active = tabIndex
        }
    }

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.headerVisible = true;
        this.tabs = TABS.tabs;
        this.setTab($state.current);
        $scope.$on("toggle-header", function () {
            self.headerVisible = !self.headerVisible;
        });
    };

});