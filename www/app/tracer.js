var log_tracer = function() {
  var obj_ref = this;

  obj_ref.timeout_obj = null;
  obj_ref.last_time_point = {term_id : null, time_value : null};

  obj_ref.get_term_id = function() {
    var id = $('#debug_id').val();
    if (isNaN(id)) {
      return null;
    }

    return id;
  };

  obj_ref.get_all_logs = function() {
    var id = obj_ref.get_term_id();
    if (id == null) {
      return;
    }

    var request = new XMLHttpRequest();

    request.onload = obj_ref.on_load;

    request.open("GET", "http://" + $('#service_addr').val() + "/terminal/" + String(id), true);
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
    // Check whether address set correctly.
  };

  obj_ref.on_load = function() {
    if (this.status != 200) {
      return;
    }

    var response_obj = JSON.parse(this.responseText);

    $('.log_entry').remove();

    response_obj.entries.forEach(function(e) {
      var message = atob(e.message);
      console.log(e.timestamp + " " + message);

      $('<tr class="log_entry"><td>' +
      e.timestamp +
      '</td><td>' +
      message +
      '</td></tr>')
        .appendTo('#log_entries_view');

      $("#logs_table").trigger("update");
    });
  };

  obj_ref.on_timeout = function() {};
};
