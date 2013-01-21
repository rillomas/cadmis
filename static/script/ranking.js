var framerate = 60;

// スケッチの状態
var SketchState = function() {
  this.sun = undefined;
  this.planets = [];
  this.transcount = 0;
  this.thetacount = 0;
  this.framecount = 0;
  this.busy = true;
}

SketchState.prototype = {
  set_sun: function(sun_) {
    this.sun = sun_;
    //for (var planet in this.planets) {
    //  console.log("planet:" + planet);
    //}
  },

  add_planet: function(planet_) {
    this.planets.push(planet_);
    //if (this.sun) {
    //  console.log("sun:" + this.sun);
    //}
  }
}

// ユーザ情報オブジェクト
var User = function(id_, rating_) {
  this.id = id_;
  this.rating = rating_;
  this.radius = undefined;
  this.x = undefined;
  this.y = undefined;
  this.r = 255;
  this.g = 0;
  this.b = 204;
  this.a = 128;
}

User.prototype = {
  type: function() {
    return "User";
  },
  get_rating: function() {
    return 3*this.rating;
  }
}

// テスト結果オブジェクト
var Test = function(id_, rating_) {
  this.id = id_;
  this.rating = rating_;
  this.radius = undefined;
  this.x = undefined;
  this.y = undefined;
  this.r = 51;
  this.g = 153;
  this.b = 255;
  this.a = 128;
}

Test.prototype = {
  type: function() {
    return "Test";
  },
  get_rating: function() {
    return 3*this.rating;
  }
}

// 座標判定
function isInCircle(x, y, obj) {
  var diff_x = obj.x - x;
  var diff_y = obj.y - y;
  if (diff_x*diff_x + diff_y*diff_y < obj.radius*obj.radius) {
    return true;
  } else {
    return false;
  }
}

function calcBias(count) {
  var x = 3 * count / framerate;
  return Math.max(0.0, Math.min(1.0, Math.log(x)+1));
}
// 結果を更新
function updateResult(target, sketch_state) {
  sketch_state.set_sun(target);
  sketch_state.planets.length = 0;
  if ("User" == target.type()) {
    sketch_state.add_planet(new Test("Q1", 10));
    sketch_state.add_planet(new Test("Q2", 20));
    sketch_state.add_planet(new Test("Q3", 30));
  } else { // "Test"
    sketch_state.add_planet(new User("Futatsugi", 10));
    sketch_state.add_planet(new User("Hori", 10));
    sketch_state.add_planet(new User("Samuel", 10));
  }
}

function textCircular(p, msg, center_x, center_y, radius, theta) {
  /* 
     for (var i=0; i<msg.length; ++i) {
     var m = msg.charAt(i);
     p.pushMatrix();

     var x = center_x + radius * Math.sin(theta);
     var y = center_y + radius * Math.cos(theta);

     p.translate(x);

     p.text

     p.popMatrix();
     }
   */
}

/**
 * ランキング画面のコントローラー
 */
function RankingController($scope, $routeParams, ranking, authenticate) {
//function RankingController($scope, $routeParams, $http) {
  
  function sketch(p) {
    p.setup = function() {
      p.size(window.innerWidth/2, window.innerHeight/1.5);
      p.frameRate(framerate);
      p.smooth();
      p.textAlign(p.CENTER, p.CENTER);
    }

    // 毎フレーム描画コールバック
    p.draw = function() {

      var should_inc_theta = true;
      
      var bias = calcBias($scope.sketch_state.transcount);

      p.background(0);
      p.fill(255, 255, 255);

      // determine center and max clock arm length
      var center_x = p.width / 2;
      var center_y = p.height / 2;

      var theta_offset = 2 * Math.PI * ($scope.sketch_state.thetacount % 3000) / 3000;

      // 遊星を描画
      (function() {
       var num = $scope.sketch_state.planets.length;
       for (var i=0; i<num; ++i) {
         var planet = $scope.sketch_state.planets[i];
         planet.x = center_x + bias * (p.width / 3) * Math.sin(theta_offset+(2 * i * Math.PI) / num);
         planet.y = center_y - bias * (p.height / 3) * Math.cos(theta_offset+(2 * i * Math.PI) / num);

         (function () {
          if (isInCircle(p.mouseX, p.mouseY, planet)) {
            planet.radius = 1.1 * planet.get_rating();
            should_inc_theta = false;
          } else {
            planet.radius = planet.get_rating();
            should_inc_theta = should_inc_theta | false;
          }

          p.stroke(255);
          p.line(planet.x, planet.y, center_x, center_y);

          p.stroke(planet.r, planet.g, planet.b);
          p.fill(planet.r, planet.g, planet.b);
          p.ellipse(planet.x, planet.y, planet.radius, planet.radius);

          p.fill(255, 255, 255);
          p.text(String(planet.id), planet.x, planet.y);
          
          p.popMatrix();
          })();
       }
       })();
      
      // 太陽を描画
      (function() {
       var sun = $scope.sketch_state.sun;
       if (sun) {
         sun.x = center_x;
         sun.y = center_y;
         sun.radius = sun.get_rating();
         
         p.stroke(sun.r, sun.g, sun.b);
         p.fill(sun.r, sun.g, sun.b);
         p.ellipse(sun.x, sun.y, sun.radius, sun.radius);
         
         p.fill(255, 255, 255);
         p.text(String(sun.id), sun.x, sun.y);
       }
      })();

      // オーバーレイ
      if ($scope.sketch_state.busy) {
        p.stroke(255, 255, 255);
        p.fill(255, 255, 255);
        p.rect(0, 0, p.width, p.height);
        p.fill(0, 0, 0);
        p.text("Loading...", p.width/2, p.height/2);
      } else {
        p.stroke(255, 255, 255, 255*(1.0-bias));
        p.fill(255, 255, 255, 255*(1.0-bias));
        p.rect(0, 0, p.width, p.height);
      }
      
      $scope.sketch_state.transcount++;
      $scope.sketch_state.framecount++;
      
      if (should_inc_theta) {
        $scope.sketch_state.thetacount++;
      }
    }

    // マウス押下時コールバック
    p.mousePressed = function() {

      var target = undefined;
           
      for (var i=0; i<$scope.sketch_state.planets.length; ++i) {
        var planet = $scope.sketch_state.planets[i];

        if (isInCircle(p.mouseX, p.mouseY, planet)) {
          target = planet;
          break;
        }
      }
      if (target) {
        
        $scope.sketch_state.busy = true;
        $scope.sketch_state.transcount = 0;
        $scope.sketch_state.planets.length = 0;

        if ("User" == target.type()) {
          ranking.getGoals(target.id, function(data) {

              for (var i in data) {
                $scope.sketch_state.planets.push(new Test(data[i].ExamId, data[i].Score/10));
              }

              $scope.sketch_state.set_sun(new User(target.id, 30));
              $scope.sketch_state.busy = false;
          });
        } else {
          ranking.getUsers(target.id, function(data) {
              
              for (var i in data) {
                $scope.sketch_state.planets.push(new User(data[i].UserId, data[i].Score/10));
              }

              $scope.sketch_state.set_sun(new Test(target.id, 30));
              $scope.sketch_state.busy = false;
          });
        }
         
        //updateResult(target, $scope.sketch_state);
      }
    }
  }
  
  $scope.sketch_state = new SketchState();
  
  authenticate.userId = 42;
  
  $scope.sketch_state.busy = true;
  $scope.sketch_state.transcount = 0;
  $scope.sketch_state.planets.length = 0;

  ranking.getGoals(authenticate.userId, function(data) {

      for (var i in data) {
        $scope.sketch_state.planets.push(new Test(data[i].ExamId, data[i].Score/10));
      }
      
      $scope.sketch_state.set_sun(new User(authenticate.userId, 30));
      $scope.sketch_state.busy = false;
  });

  var canvas = document.getElementById("ranking-canvas");
  var p = new Processing(canvas, sketch);
  // p.exit(); to detach it
}
