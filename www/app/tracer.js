var tracer = function() {
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

    response_obj.entries.forEach(function(e) {
      console.log(e.timestamp + " " + atob(e.message));
    });
  };

  request.open("GET", "http://" + $('#service_addr').val() + "/terminal/" + String(id), true);
  request.send();
};
