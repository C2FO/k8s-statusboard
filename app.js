
// Initially get the contexts and add divs for them to the dom.
$.get("/api/contexts", function(data, status){
  data.forEach(function(context){
    $("#context-container").append('<div id="context-' + context + '" ></div>');
  });
});

// Event source to 
var es = new EventSource("/events/");
es.addEventListener("pod-status", function(e){
  var obj = JSON.parse(e.data);

  var container = $("div#context-" + obj.context).empty();
  obj.pods.forEach(function(pod){
    console.log(pod);
    container.append("<p>" + pod.metadata.name + "</p>");
  });
});