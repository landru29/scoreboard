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
    this.resolvers = {};
    this.uuid = null;

    this.autoReconnect = false;
    this.maxReconnectionAttempts = 10;

    if (!window.WebSocket) {
        console.error("Websockets does not Work on this Browser ! Use another browser like Firefox or Chrome.");
    }

    this.generateUuid = function () {
        return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function(c) {
            var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8);
            return v.toString(16);
        });
    }

    /**
     * Create a websocket connection
     * @return {Promise}
     */
    this.createConnection = function (name, url) {
        var self = this;
        return new Promise(function (resolve, reject) {
            url = /^ws:\/\//.test(url) ? url : "ws://" + window.location.host + url;
            self.connections[name] = new WebSocket(url);
            self.attempts[name] = !_.isUndefined(self.attempts[name]) ? self.attempts[name] : self.maxReconnectionAttempts;
            self.connections[name].onopen = function () {
                resolve(self.connections[name]);
            }
            self.connections[name].onerror = function () {
                reject();
            }
            self.connections[name].onclose = function (evt) {
                console.warn("[Websocket]", "connection closed");
                if (self.attempts[name]-- >= 0) {
                    console.log("[Websocket]", "reconnecting");
                    self.createConnection (name, url);
                    return;
                }
                console.error("[Websocket]", "Fail to connect to host");
                delete self.connections[name];
            };
            self.connections[name].onmessage = function (evt) {
                var response = evt.data;
                try {
                    response = JSON.parse(evt.data);
                    response.data = JSON.parse(response.data);
                } catch (e) {}
                var requestId = _.get(response, "requestId");
                if (requestId && self.resolvers[requestId]) {
                    self.resolvers[requestId](response);
                    delete self.resolvers[requestId];
                }
                switch (_.get(response, "status")) {
                    case "sync":
                        self.offset = {
                            clientToServer: response.data.offset,
                            serverToClient: new Date().getTime() - response.timestamp,
                            clientToClient: new Date().getTime() - response.data.clientTimestamp
                        }
                        break;
                    case "whoami":
                        self.uuid = _.get(response, "data.uuid");
                        break;
                    default:
                }
                if (self.$rootScope) {
                    self.$rootScope.$broadcast("ws-message", response);
                    _.forEach(self.messageTransporters, function(cb, key) {
                        cb(response, key);
                    });
                }
            };
        });
    }

    this.send = function (name, command, data) {
        var uuid = this.generateUuid();
        var commandToSend = JSON.stringify({
            name: command,
            data: _.isObject(data) ? JSON.stringify(data) : data,
            timestamp: new Date().getTime(),
            requestId: uuid
        });
        
        if (!this.connections[name]) {
            console.error("[Websockets]", "Connection", name, "does not exists");
            return;
        }
        //console.info("[Websockets]", "sending", commandToSend);
        return new Promise(function (resolve) {
            self.connections[name].send(commandToSend);
            self.resolvers[uuid] = resolve;
        });
    }

    var Ws = function Ws() {
    };

    Ws.prototype.getUuid = function () {
        return self.uuid;
    };

    Ws.prototype.createConnection = function (name, url) {
        self.createConnection(name, url);
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
        self.send(name, command, data);
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