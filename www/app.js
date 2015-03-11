var fire_timer = function() {
  alert('it is time!');
};

var start_trace = function() {
  // setInterval(fire_timer, 6000);
  var id = $('#debug_id').val();
  if (isNaN(id)) {
    alert('Not a number!');
  } else {
    alert(+id);
  }
};

var submit_debug_id_by_enter = function(event) {
  if (event.which == 13) {
    alert('Submit fired!');
    $('#start_trace_btn').click();
  }
};

var on_loaded = function() {
  $('#start_trace_btn').on('click', start_trace);
  $('#debug_id').keyup(submit_debug_id_by_enter);
};

$(document).ready(on_loaded);
