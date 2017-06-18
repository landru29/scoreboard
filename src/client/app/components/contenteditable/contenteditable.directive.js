angular.module("scoreboard").directive("contenteditable", function directive () {
    "use strict";
    return {
            restrict: "A",
            require: "ngModel",
            link: function (scope, element, attrs, ngModel) {
                function read() {
                    // view -> model
                    var html = element.html();
                    html = html.replace(/&nbsp;/g, "\u00a0");
                    if (attrs.type === "number") {
                        var display = "" + (parseInt(html, 10) || 0);
                        ngModel.$setViewValue(parseInt(html, 10));
                        element.html(display);
                    } else {
                        ngModel.$setViewValue(html);
                    }
                }
                // model -> view
                ngModel.$render = function() {
                    var display;
                    if (attrs.type === "number") {
                        display = parseInt("" + ngModel.$viewValue, 10) || "0";
                    } else {
                        display = _.isUndefined(ngModel.$viewValue) ? "" : "" + ngModel.$viewValue | "";
                    }
                    element.html(display);
                };

                element.bind("blur", function() {
                    scope.$apply(read);
                });
                element.bind("keydown keypress", function (event) {
                    if(event.which === 13) {
                        this.blur();
                        event.preventDefault();
                    }
                });
            }
        };
});