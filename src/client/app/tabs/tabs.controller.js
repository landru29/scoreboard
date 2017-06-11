angular.module("scoreboard").controller("TabsCtrl", function TabsCtrl ($scope, TABS, $state) {
    "use strict";

    this.go = function (state) {
        return $state.go(state);
    }

    /**
     * Initialization of the controller
     */
    this.$onInit = function () {
        this.tabs = TABS.tabs;
        var tabName;
        var matcher = $state.current.name.match(/^(main\.tabs\.[\w]*)/);
        if (matcher) {
            tabName = matcher[1];
        }
        var tabIndex = _.findIndex(TABS.tabs, { state: tabName });
        if (tabIndex > -1) {
            this.active = tabIndex
        }
    };

});