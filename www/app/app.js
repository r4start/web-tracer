var ids_list_fetcher = new ids_fetcher();
var log_fetcher = new log_tracer();

var auto_updater_observers = [];

var submit_debug_id_by_enter = function(event) {
  if (event.which != 13) {
    return;
  }

  log_fetcher.get_all_logs();
};

var get_ids = function (event) {
  if (ids_list_fetcher.fetch_started) {
    return;
  }

  if (event.which != 13 && $(event.target).is("#service_addr")) {
    return;
  }

  if (ids_list_fetcher.timeout_obj == null) {
    ids_list_fetcher.fetch_ids();
  }
};

var auto_updater_toggled = function() {
  var state = this.checked;
  auto_updater_observers.forEach(function(elem) {
    elem(state);
  });
};

var on_loaded = function() {
  auto_updater_observers.push(log_fetcher.auto_update_observer);

  var dbg_id_field = $('#debug_id');
  dbg_id_field.keyup(submit_debug_id_by_enter);

  dbg_id_field.focusin(function (event) {
    dbg_id_field.val("");
  });

  dbg_id_field.focus(get_ids);

  var service_field = $('#service_addr');
  service_field.keyup(get_ids);
  service_field.val(location.host);

  $('#autoupdate').click(auto_updater_toggled);
};

$(document).ready(on_loaded);
