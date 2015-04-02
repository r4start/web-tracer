var ids_fetcher = new IdsFetcher();

var submit_debug_id_by_enter = function(event) {
  if (event.which != 13) {
    return;
  }

  tracer();
};

var get_ids = function (event) {
  if (ids_fetcher.fetch_started) {
    return;
  }

  if (event.which != 13 && $(event.target).is("#service_addr")) {
    return;
  }

  ids_fetcher.fetch_ids();
};

var on_loaded = function() {
  var dbg_id_field = $('#debug_id');
  dbg_id_field.keyup(submit_debug_id_by_enter);

  dbg_id_field.focusin(function (event) {
    dbg_id_field.val("");
  });

  dbg_id_field.focus(get_ids);

  var service_field = $('#service_addr');
  service_field.keyup(get_ids);
  service_field.val(location.host);

  $('#logs_table').tablesorter();
};

$(document).ready(on_loaded);
