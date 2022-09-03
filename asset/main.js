

var startRecordingTime = -1
var inter = -1

$(document).ready(function () { 
  $( "#dialog" ).dialog(
    { autoOpen: false,
    draggable: false,
    height: 400
  });
  
  $( "#autocomplete" ).autocomplete();

  $.ajax({
    url: 'getstate',
    type: 'get',
    success : function(data) {
      var actual = JSON.parse(atob(data))
      console.log(actual.isRecording)
      console.log(actual.serials)
      $( "#autocomplete" ).autocomplete( "option", "source", actual.serials );
      if(actual.isRecording) {
        displayRecording(actual.batterySerial)
        startRecordingTime = actual.startRecordingTime
      }
    },
  });
});
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
    // open dialog:
    $( "#dialog" ).dialog( "open")
	}
	else{
    endRecording()
	}
});
/*
$(function() {
  $( "#autocomplete" ).autocomplete({
     source: [
        { label: "a", value: "a" },
        { label: "b", value: "b" }
     ]
  });
});
*/
$('#submit-serial').click(function () {
  $.ajax({
    url: 'newrec',
    type: 'post',
    dataType: 'html',
    data : { serial: $('#autocomplete').val()},
    success : function(data) {
      // change button style:
      startRecordingTime = Date.now() / 1000
      displayRecording($('#autocomplete').val())
      $( "#dialog" ).dialog( "close")
    },
  });
});

function endRecording() {
  $.ajax({
    url: 'endrec',
    type: 'get',
    dataType: 'html',
    success : function(data) {
      // change button style:
      $('#recButton').removeClass("Rec");
      $('#recButton').addClass("notRec");
      $('#rec-stat').text("Start Recording");
      $('#recCard').removeClass("Rec");
      clearInterval(inter);
    },
  });
};

function displayRecording(serial) {
  $('#recButton').removeClass("notRec");
      $('#recButton').addClass("Rec");
      $('#rec-stat').text("00:00:00 \nRecording battery: " + serial);
      $('#recCard').addClass("Rec");
      function foo () { 
      $('#rec-stat').text(formatTime() + "\nRecording battery: " + serial);
     }
     inter = setInterval(foo, 1000);
}

function formatTime() {
  d = (Date.now()/1000) - startRecordingTime
  if(d <= 0){
    return '00:00:00'
  }else{
   let h = Math.floor(d / 3600);
   let m = Math.floor(d % 3600 / 60);
   let s = Math.floor(d % 3600 % 60);

   let hDisplay = h <= 9 ? '0'+ h+':' : h+ ":";
   let mDisplay = m <= 9 ? '0'+ m+':' : m+ ":";
   let sDisplay = s <= 9 ? '0'+ s : s;


   return hDisplay + mDisplay + sDisplay; 
  }
}

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