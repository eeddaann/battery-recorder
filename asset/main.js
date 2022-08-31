
// generate sample data and put it in array
var data = [];

// create the graph
var g = new Dygraph(document.getElementById("div_g"), data,
                    {
                      drawPoints: true,
                      showRoller: false,
                      series: {
                        'Temperature': {
                          
                          axis: 'y2'
                        },
                        'Voltage': {
                          axis: 'y'
                        },
                      },
                      axes: {
                        y: {
                          independentTicks: true
                        },
                        y2: {
                          // set axis-related properties here
                          independentTicks: true
                        }
                      },
                      labels: ['Time', 'Temperature', 'Voltage']
                    });

$('#recButton').addClass("notRec");

$('#recCard').click(function(){
	if($('#recButton').hasClass('notRec')){
		$('#recButton').removeClass("notRec");
		$('#recButton').addClass("Rec");
    $('#rec-stat').text("Recording...");
    $('#recCard').addClass("Rec");
	}
	else{
		$('#recButton').removeClass("Rec");
		$('#recButton').addClass("notRec");
    $('#rec-stat').text("Start Recording");
    $('#recCard').removeClass("Rec");
	}
});	


var socket = io(); // create socket
socket.on('temp', function(msg){
  var x = new Date();  // current time
  tmp = msg.split(',')
  var temperature = parseFloat(tmp[0]);
  var voltage = parseFloat(tmp[1]);
  if(data.length > 105) { // delete random element
    const random = Math.floor(Math.random() * 100)+1; // sample one of the first 100 points, keep first point
    data.splice(random, 1)[0];
  }
  data.push([x, temperature, voltage]);
  g.updateOptions( { 'file': data } );
  $('#volt_cnt').text(Math.round(voltage * 100) / 100 + "V");
  $('#temp_cnt').text(Math.round(temperature * 100) / 100 + "°C");
});