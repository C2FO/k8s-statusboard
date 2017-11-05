function Metrics(pods) {
  this.running = 0;
  this.notRunning = 0;
  this.succeeded = 0;
  if (arguments.length > 0) {
    for (var i=0; i < pods.length; i ++){
      var phase = pods[i].status.phase;
      if (phase == "Succeeded"){
        this.succeeded += 1;
      } else if (phase == "Running"){
        this.running += 1;
      } else {
        this.notRunning += 1;
      }
    }
  }
}

function Context(name) {
  var self = this;
  this.name = name;
  this.metrics = new Metrics();
  this.update = function(pods) {
    self.metrics = new Metrics(pods);
  };
}

// Begin angular App
angular.module('k8sStatusApp', [])
.controller('mainController', function($scope, $http) {
    // create a message to display in our view
    $scope.title = 'K8S StatusBoard';
    $scope.contexts = [];

    // Load the initial contexts
    $http.get("/api/contexts")
      .success(function(data){
        data.forEach(function(contextName){
          $scope.contexts.push(new Context(contextName));
        });
      });

    // Event source to 
    var es = new EventSource("/events/");
    es.addEventListener("pod-status", function(e){
      var obj = JSON.parse(e.data);
      for(var i = 0; i < $scope.contexts.length; i++){
        if($scope.contexts[i].name == obj.context) {
          $scope.contexts[i].update(obj.pods);
          $scope.$apply();
        }
      }
    });
})
.directive('contextMetrics', function(){
  return {
    restrict: 'E',
    replace: true,
    scope: {
      name: '=',
      metrics: '='
    },
    templateUrl: 'context-metrics.html',
  };
});