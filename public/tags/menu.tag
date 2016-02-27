<menu>
   <div class="navbar-fixed">
      <nav>
         <div class="nav-wrapper">
            <ul class="left">
               <li each={ options } class={ active: parent.current === url }><a href={"#"+url}>{ label } </a></li>
               <li><span class="brand">ipc-schedule</span></li>
            </ul>
         </div>
      </nav>
   </div>

   var self = this
   self.current = false

   this.on('mount', function() {
      self.current = location.hash
      self.update()
   })

   riot.route.create()(function(url) {
      console.log("New route is",url)
      self.current = url
      self.update()
   })

   this.options = [
      { label: 'Home', url: '' },
      { label: 'Groups', url: 'groups' },
      { label: 'Upload', url: 'upload' },
      { label: 'API', url: 'usage' },
   ]
</menu>
