var start_trace = function() {
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

var get_ids = function(event) {
  if (event.which != 13) {
    return
  }

  ids_fetcher();
};

var on_loaded = function() {
  $('#start_trace_btn').on('click', start_trace);
  $('#debug_id').keyup(submit_debug_id_by_enter);
  $('#service_addr').keyup(get_ids);
};

$(document).ready(on_loaded);
