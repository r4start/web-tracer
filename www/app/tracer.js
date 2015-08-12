var log_tracer = function() {
  var obj_ref = this;

  obj_ref.timeout_obj = null;
  obj_ref.auto_update_state = {term_id : null, time_value : null, service_address : null};

  obj_ref.get_term_id = function() {
    var id = $('#debug_id').val();
    if (isNaN(id)) {
      return null;
    }

    return id;
  };

  obj_ref.get_service_address = function() {
    var addr = $('#service_addr').val();
    return addr;
  };

  obj_ref.get_all_logs = function() {
    var id = obj_ref.get_term_id();
    if (id == null) {
      return;
    }

    var request = new XMLHttpRequest();

    request.onload = obj_ref.on_load;

    request.open(
          "GET",
          "http://" + obj_ref.get_service_address() + "/terminal/" + String(id),
          true);
    request.send();
  };

  obj_ref.auto_update_observer = function(state) {
    if (obj_ref.timeout_obj != null) {
      clearTimeout(obj_ref.timeout_obj);
      obj_ref.timeout_obj = null;
    }

    if (state === false) {
      return;
    }

    var term_id = obj_ref.get_term_id();
    if (term_id == null) {
      return;
    }

    var addr = obj_ref.get_service_address();
    if (addr == null || addr === "") {
      return;
    }

    obj_ref.auto_update_state.term_id = term_id;
    obj_ref.auto_update_state.service_address = addr;

    setTimeout(obj_ref.auto_update_on_timeout);
  };

  obj_ref.on_load = function() {
    if (this.status != 200) {
      return;
    }

    var response_obj = JSON.parse(this.responseText);

    $('.log_entry').remove();

    var latest_time = null;

    response_obj.entries.forEach(function(e) {
      var message = atob(e.message);
      console.log(e.timestamp + " " + message);

      if (latest_time == null) {
        latest_time = e.timestamp;
      } else if (e.timestamp > latest_time) {
        latest_time = e.timestamp;
      }

      $('<tr class="log_entry"><td>' +
      e.timestamp +
      '</td><td>' +
      message +
      '</td></tr>')
        .appendTo('#log_entries_view');

      $("#logs_table").trigger("update");
    });

    console.log('Latest timestamp ' + latest_time);
    if (latest_time > obj_ref.auto_update_state.time_value) {
      obj_ref.auto_update_state.time_value = latest_time;
    }
  };

  obj_ref.auto_update_on_timeout = function() {
    console.log("Update data from server. " +
                obj_ref.auto_update_state.term_id + " " +
                obj_ref.auto_update_state.service_address + " " +
                obj_ref.auto_update_state.latest_time);
  };
};
