var submit_debug_id_by_enter = function(event) {
  if (event.which != 13) {
    return;
  }

  $('#start_trace_btn').click();
};

var get_ids = function(event) {
  if (event.which != 13) {
    return
  }

  ids_fetcher();
};

var on_loaded = function() {
  $('#start_trace_btn').on('click', tracer);
  $('#debug_id').keyup(submit_debug_id_by_enter);
  $('#service_addr').keyup(get_ids);
  $('#logs_table').tablesorter();
};

$(document).ready(on_loaded);
