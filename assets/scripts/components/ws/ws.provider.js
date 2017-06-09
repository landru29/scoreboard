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
    this.attempts = {};
    this.$rootScope = null

    this.autoReconnect = false;
    this.maxReconnectionAttempts = 10;

    if (!window.WebSocket) {
        console.error("Websockets does not Work on this Browser ! Use another browser like Firefox or Chrome.");
    }

    this.createConnection = function (name, url) {
        var self = this;
        url = /^ws:\/\//.test(url) ? url : "ws://" + window.location.host + url;
        this.connections[name] = new WebSocket(url);
        this.attempts[name] = this.attempts[name] ? this.attempts[name] : this.maxReconnectionAttempts;
        this.connections[name].onclose = function (evt) {
            console.warn("[Websocket]", "connection closed");
            if (self.attempts[name]-- >= 0) {
                console.log("[Websocket]", "reconnecting");
                self.createConnection (name, url);
                return;
            }
            console.error("[Websocket]", "Fail to connect to host");
            delete self.connections[name];
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

    Ws.prototype.send = function (name, command, data) {
        var commandToSend = JSON.stringify({
            name: command,
            data: _.isObject(data) ? JSON.stringify(data) : data
        });
        
        if (!self.connections[name]) {
            console.error("[Websockets]", "Connection", name, "does not exists");
            return;
        }
        console.info("[Websockets]", "sending", commandToSend);
        self.connections[name].send(commandToSend);
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