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
        // Now when the phase is running, we still need to check the array of
        // pod.status.containerStatuses
        var isRunning = true;
        var containerStatuses = pods[i].status.containerStatuses;
        for (var j=0; j < containerStatuses.length; j++) {
          if (!containerStatuses[j].state.hasOwnProperty('running')) {
            // at least one of the containers is not Running.
            isRunning = false;
          }
        }
        if (isRunning) {
          this.running += 1;
        } else {
          this.notRunning += 1;
          console.log("Pod " + pods[i].metadata.name + " is not running");
        }
      } else {
        this.notRunning += 1;
        console.log("Pod " + pods[i].metadata.name + " is not running");
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
          $scope.$apply(); // Import to get re-renders.
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
    controller: [ "$scope", function($scope){
      $scope.contextClass = function() {
        if($scope.metrics.notRunning == 0) {
          return 'green';
        } else {
          return 'red';
        }
      };
    }],
    templateUrl: 'context-metrics.html'
  };
});