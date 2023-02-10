const watch = document.querySelector("#stopwatch");
      let millisecound = 0;
      let timer;

      function timeStart(){
        document.getElementById("start").style.display = "none";
        document.getElementById("pause").style.display = "inline";
        watch.style.color = "#0f62fe";
        clearInterval(timer);
        timer = setInterval(() => {
          millisecound += 10;

          let dateTimer = new Date(millisecound);

          watch.innerHTML =
          ('0'+dateTimer.getUTCHours()).slice(-2) + ':' +
          ('0'+dateTimer.getUTCMinutes()).slice(-2) + ':' +
          ('0'+dateTimer.getUTCSeconds()).slice(-2) + ':' +
          ('0'+dateTimer.getUTCMilliseconds()).slice(-3,-1);
        }, 10);
      }


      function timePaused() {
        watch.style.color = "red";
        clearInterval(timer);
        document.getElementById("start").style.display = "inline";
        document.getElementById("pause").style.display = "none";
      }

      document.addEventListener('click', (e) => {
        const el = e.target;

        if(el.id === 'start') timeStart();
        if(el.id === 'pause') timePaused();
      })
