var log_tracer = function() {
  var obj_ref = this;

  obj_ref.get_all_logs = function() {
    var id = $('#debug_id').val();
    if (isNaN(id)) {
      alert('Not a number!');
      return;
    }

    var request = new XMLHttpRequest();

    request.onload = function() {
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

    request.open("GET", "http://" + $('#service_addr').val() + "/terminal/" + String(id), true);
    request.send();
  };
};
