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
    console.log(response_obj.entries);
  };

  request.open("GET", "http://" + $('#service_addr').val() + "/terminal/" + String(id), true);
  request.send();
};
