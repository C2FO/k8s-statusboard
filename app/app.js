
// var app = angular.module('k8sStatusApp');
// app.controller('mainController', function($scope) {
//     // create a message to display in our view
//     $scope.title = 'K8S StatusBoard';
// });

// Initially get the contexts and add divs for them to the dom.
$.get("/api/contexts", function(data, status){
  data.forEach(function(context){
    $("#context-container").append('<div id="context-' + context + '" ></div>');
  });
});

// getRuningFailed returns the number of running pods and the number of failed
// pods
function getPodMetrics(pods) {
  var metrics = {
    running: 0,
    succeeded: 0,
    notRunning: 0,
  };
  pods.forEach(function(pod){
    if (pod.status.phase == "Succeeded"){
      metrics.succeeded += 1;
    } else if (pod.status.phase == "Running") {
      metrics.running += 1;
    } else {
      metrics.notRunning += 1;
    }
  });
  return metrics;
}

// Event source to 
var es = new EventSource("/events/");
es.addEventListener("pod-status", function(e){
  var obj = JSON.parse(e.data);
  var container = $("div#context-" + obj.context).empty();
  var metrics = getPodMetrics(obj.pods);
  container.first().html('<div class="col s4"><h4>' + obj.context + "</h4>" + "<p>Running: " + metrics.running + " Not Running: " + metrics.notRunning + "</p></div>");
});