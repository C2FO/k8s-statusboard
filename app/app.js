function PodMetrics(pods) {
  this.running = 0;
  this.notRunning = 0;
  this.succeeded = 0;

  if (typeof pods == 'undefined' || pods === null) {
    return;
  }

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

function JobMetrics(jobs) {
  this.desired = 0;
  this.succeeded = 0;

  if (typeof jobs == 'undefined' || jobs === null) {
    return;
  }

  if (arguments.length > 0) {
    for (var i=0; i < jobs.length; i++){
      var status = jobs[i].status;
      var spec = jobs[i].spec;
      
      if (spec.hasOwnProperty('completions')) {
        this.desired += spec.completions;
      }

      if (status.hasOwnProperty('succeeded')) {
        this.succeeded += status.succeeded;
      }
    }
  }
}

function ServiceMetrics(services) {
  this.loadBalancer = 0;
  this.externalName = 0;
  this.clusterIP = 0;
  this.nodePort = 0;

  if (typeof services == 'undefined' || services === null) {
    return;
  }

  if (arguments.length > 0) {
    for (var i=0; i < services.length; i++){
      var spec = services[i].spec;

      if (spec.hasOwnProperty('type')) {
        var type = spec.type;
        if (type === 'NodePort') {
          this.nodePort += 1;
        } else if (type === 'ClusterIP') {
          this.clusterIP += 1;
        } else if (type === 'ExternalName') {
          this.externalName += 1;
        } else if (type === 'LoadBalancer') {
          this.loadBalancer += 1;
        } else {
          console.log(type);
        }
      }
    }
  }
}

function Context(name) {
  var self = this;
  this.name = name;
  this.pod_metrics = new PodMetrics();
  this.job_metrics = new JobMetrics();
  this.service_metrics = new ServiceMetrics();

  this.updatePods = function(pods) {
    self.pod_metrics = new PodMetrics(pods);
  };
  this.updateJobs = function(jobs) {
    self.job_metrics = new JobMetrics(jobs);
  };
  this.updateServices = function(services) {
    self.service_metrics = new ServiceMetrics(services);
  };
}

// Begin angular App
angular.module('k8sStatusApp', [])
.controller('mainController', function($scope, $http) {
    // create a message to display in our view
    $scope.title = 'K8S StatusBoard';
    $scope.contexts = [];

    // Load the initial contexts
    $http.get("./api/contexts")
      .success(function(data){
        data.forEach(function(contextName){
          $scope.contexts.push(new Context(contextName));
        });
      });

    // Event source to 
    var es = new EventSource("./events/");
    es.addEventListener("pod-status", function(e){
      var obj = JSON.parse(e.data);
      for(var i = 0; i < $scope.contexts.length; i++){
        if($scope.contexts[i].name == obj.context) {
          $scope.contexts[i].updatePods(obj.pods);
          $scope.$apply(); // Important to get re-renders.
        }
      }
    });
    
    es.addEventListener("job-status", function(e){
      var obj = JSON.parse(e.data);
      for(var i = 0; i < $scope.contexts.length; i++){
        if($scope.contexts[i].name == obj.context) {
          $scope.contexts[i].updateJobs(obj.jobs);
          $scope.$apply(); // Important to get re-renders.
        }
      }
    });

    es.addEventListener("services-status", function(e){
      var obj = JSON.parse(e.data);
      for(var i = 0; i < $scope.contexts.length; i++){
        if($scope.contexts[i].name == obj.context) {
          $scope.contexts[i].updateServices(obj.services);
          $scope.$apply(); // Important to get re-renders.
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
      podMetrics: '=',
      jobMetrics: '=',
      serviceMetrics: '='
    },
    controller: [ "$scope", function($scope){
      $scope.contextClass = function() {
        if($scope.podMetrics.notRunning == 0) {
          if($scope.jobMetrics.desired <= $scope.jobMetrics.succeeded) {
            return 'green';
          }
        }
        return 'red';
      };
    }],
    templateUrl: 'context-metrics.html'
  };
});