
// generate sample data and put it in array
var data = [];
var t = new Date();
for (var i = 10; i >= 0; i--) {
  var x = new Date(t.getTime() - i * 1000);
  data.push([x, Math.random()]);
}

// create the graph
var g = new Dygraph(document.getElementById("div_g"), data,
                    {
                      drawPoints: true,
                      showRoller: true,
                      valueRange: [0.0, 100],
                      labels: ['Time', 'Random']
                    });


var socket = io(); // create socket
socket.on('temp', function(msg){
  $('#temp').text(msg);
  var x = new Date();  // current time
  var y = msg;
  if(data.length > 105) { // delete random element
    const random = Math.floor(Math.random() * 100)+1; // sample one of the first 100 points, keep first point
    data.splice(random, 1)[0];
  }
  data.push([x, y]);

  g.updateOptions( { 'file': data } );
});