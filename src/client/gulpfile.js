// Requis
var gulp = require("gulp");

// Include plugins
var plugins = require("gulp-load-plugins")();
var gutil = require("gulp-util");

console.log(Object.keys(plugins));

var destination = "../../dist/client";
var dev = "./dev";

function string_src(filename, string) {
  var src = require("stream").Readable({ objectMode: true })
  src._read = function () {
    this.push(new gutil.File({
      cwd: "",
      base: "",
      path: filename,
      contents: new Buffer(string)
    }))
    this.push(null)
  }
  return src
}


gulp.task("clean-dist", function () {
    return gulp.src(destination, { read: false })
        .pipe(plugins.clean({ force: true }));
});

gulp.task("clean-dev", function () {
    return gulp.src(dev, { read: false })
        .pipe(plugins.clean({ force: true }));
});

gulp.task("copy-assets", function () {
    return gulp
    .src(["assets/**/*"])
    .pipe(gulp.dest(destination))
    .pipe(gulp.dest(dev + "/app"));
});

gulp.task("copy-html", function () {
    return gulp
    .src(["app/**/*.html"])
    .pipe(gulp.dest(dev + "/app"));
});

gulp.task("babel", function () {
    return gulp
    .src(["app/**/*.js", "!./app/templates.js"])
    .pipe(plugins.babel({
        presets: ["es2015"]
    }))
    .pipe(gulp.dest(destination + "/app"))
    .pipe(gulp.dest(dev + "/app"));
});

gulp.task("copy-vendor", function () {
    return gulp
    .src(["vendor/**/*.js"])
    .pipe(gulp.dest(destination + "/vendor"))
    .pipe(gulp.dest(dev + "/vendor"));
});

gulp.task("copy-fonts", function () {
    return gulp
    .src(["bower_components/bootstrap/fonts/*.{ttf,woff,woff2,eot,otf,svg}", "bower_components/font-awesome/fonts/*.{ttf,woff,woff2,eot,otf,svg}"])
    .pipe(gulp.dest(destination + "/fonts"))
    .pipe(gulp.dest(dev + "/fonts"));
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
    .pipe(plugins.rename(dev + "/index.html"))
    .pipe(gulp.dest("./"))
    .pipe(plugins.livereload());
});

gulp.task("inject-app-js", ["inject-dev-js"], function () {
    return gulp.src(dev + "/dev.html")
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
    .pipe(gulp.dest(dev + "/style"));
});

gulp.task("template-dev", function () {
  return string_src("templates.js", "angular.module(\"scoreboard.templates\", []);")
    .pipe(gulp.dest(dev + "/app"));
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
        .pipe(gulp.dest(destination + "/bower_components"))
        .pipe(gulp.dest(dev + "/bower_components"));
});  

gulp.task("watch", ["inject-dev-js", "less"], function () {
    plugins.livereload.listen();
    gulp.watch("./app/**/*.{js,less}", ["inject-dev-js", "less", "babel"])
});

gulp.task("build", [
    "inject-app-js",
    "copy-assets",
    "copy-bower-deps",
    "babel",
    "template-dev",
    "copy-html",
    "copy-vendor",
    "copy-fonts",
    "less",
    "templates"
]);

gulp.task("default", function(done) {
    plugins.runSequence("clean-dev", "clean-dist", "build");
});

gulp.task("serve", ["build", "watch"]);