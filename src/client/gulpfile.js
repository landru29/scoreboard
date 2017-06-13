// Requis
var gulp = require("gulp");

// Include plugins
var plugins = require("gulp-load-plugins")();

console.log(Object.keys(plugins));

var destination = "../../dist/client";


gulp.task("clean", function () {
    return gulp.src(destination, { read: false })
        .pipe(plugins.clean({ force: true }));
});

gulp.task("copy-assets", function () {
    return gulp
    .src(["assets/**/*"])
    .pipe(gulp.dest(destination));
});

gulp.task("copy-src", function () {
    return gulp
    .src(["app/**/*.js", "!./app/templates.js"])
    .pipe(gulp.dest(destination + "/app"));
});

gulp.task("copy-vendor", function () {
    return gulp
    .src(["vendor/**/*.js"])
    .pipe(gulp.dest(destination + "/vendor"));
});

gulp.task("copy-fonts", function () {
    return gulp
    .src(["node_modules/bootstrap/fonts/*.{ttf,woff,woff2,eot,otf,svg}", "node_modules/font-awesome/fonts/*.{ttf,woff,woff2,eot,otf,svg}"])
    .pipe(gulp.dest(destination + "/fonts"))
    .pipe(gulp.dest("./fonts"));
});

gulp.task("inject-dev-js", function () {
    return gulp.src("./index.html")
    .pipe(plugins.inject(gulp.src(
        [
            "./app/app.js",
            "./app/**/*.js",
            "./vendor/**/*.js",
            "!./app/templates.js"
        ],
        {
            read: false
        }
    ), {
        relative: true,
        name: "app"
    }
    ))
    .pipe(plugins.wiredep({}))
    .pipe(plugins.rename("dev.html"))
    .pipe(gulp.dest("./"))
    .pipe(plugins.livereload());
});

gulp.task("inject-app-js", ["inject-dev-js"], function () {
    return gulp.src("./dev.html")
        .pipe(plugins.rename("index.html"))
        .pipe(gulp.dest(destination));
});

gulp.task("less", function () {
  return gulp.src("./app/**/*.less")
    .pipe(plugins.wiredep({}))
    .pipe(plugins.less({ path: []}))
    .pipe(plugins.concat("main.css"))
    .pipe(plugins.cleanCss({ compatibility: "ie8" }))
    .pipe(gulp.dest(destination + "/style"))
    .pipe(gulp.dest("./style"));
});

gulp.task("templates", function() {
  gulp.src("{*/,**/}*.html", { cwd: "app" })
    .pipe(plugins.htmlmin({ collapseWhitespace: true }))
    .pipe(plugins.ngTemplates({
            filename: "templates.js",
            module: "scoreboard.templates",
            path: function (path, base) {
                return path.replace(base, 'app/');
            }
        }))
    .pipe(gulp.dest(destination + "/app"));
});

gulp.task("copy-bower-deps", function() {
    return gulp
        .src(["bower_components/**/*.{js,css}"])
        .pipe(gulp.dest(destination + "/bower_components"));
});  

gulp.task("watch", ["inject-dev-js", "less"], function () {
    plugins.livereload.listen();
    gulp.watch("./app/**/*.{js,less}", ["inject-dev-js", "less"])
});

gulp.task("build", [
    "inject-app-js",
    "copy-assets",
    "copy-bower-deps",
    "copy-src",
    "copy-vendor",
    "copy-fonts",
    "less",
    "templates"
]);

gulp.task("default", function(done) {
    plugins.runSequence("clean", "build");
});

gulp.task("serve", ["copy-fonts", "watch"]);