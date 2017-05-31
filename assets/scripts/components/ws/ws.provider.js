/**
 * @ngdoc service
 * @name scoreboard.service:WsProvider
 * @description
 * <p>Description</p>
 */
angular.module("scoreboard").provider("Ws", function Ws () {
    "use strict";
    var self = this;

    this.connections = {};
    this.messageTransporters = {};
    this.$rootScope = null

    this.autoReconnect = false;
    this.maxReconnectionAttempts = 10;

    if (!windows.WebSocket) {
        console.error("Websockets does not Work on this Browser ! Use another browser like Firefox or Chrome.");
    }

    this.CreateConnection = function (name, url) {
        this.connections[name] = new WebSocket(url);
        this.connections[name].onclose = function (evt) {
            console.warn("[Websocket]", "connection closed");
            if (_ws.maxReconnectionAttempts-- >= 0) {
                console.log("[Websocket]", "reconnecting");
                self.CreateConnection (name, url);
            }
        };
        this.connections[name].onmessage = function (evt) {
            var data = evt.data;
            try {
                data = JSON.parse(evt.data);
            } catch (e) {}
            if (self.$rootScope) {
                self.$rootScope.$broadcast("ws-message", data);
                _.forEach(self.messageTransporters, function(cb, key) {
                    cb(data, key);
                });
            }
        };
    }

    var Ws = function Ws() {
    };

    Ws.prototype.AddTransporter = function (name, cb) {
        self.messageTransporters[name] = cb;
    };

    Ws.prototype.removeTransporter = function (name) {
        delete self.messageTransporters[name];
    };

    Ws.prototype.registerScope = function(scope) {
        self.$rootScope = scope;
    }

    Ws.prototype.send = function (name, data) {
        data = _.isObject(data) ? JSON.stringify(data) : data;
        if (!self.connections[name]) {
            console.error("[Websockets]", "Connection", name, "does not exists");
            return;
        }
        self.connections[name].send(data);
    };

   /**
    * @ngdoc service
    * @name scoreboard.service:Ws
    * @description
    * <p>Description</p>
    */
    this.$get = function WsFn ($rootScope) {
        return new Ws();
    };
});