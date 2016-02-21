<menu>
   <div class="navbar-fixed">
      <nav>
         <div class="nav-wrapper">
            <ul class="left">
               <li each={ options } class={ active: url }><a href={"#"+url}>{ label }</a></li>
            </ul>
         </div>
      </nav>
   </div>

   riot.route(function(url) {
      console.log("New route is",url)
      this.current = url
   })

   this.options = [
      { label: 'Home', url: '' },
      { label: 'Groups', url: 'groups' },
      { label: 'Upload', url: 'upload' },
      { label: 'API', url: 'usage' },
   ]

   isActive() {
      console.log("current",this.parent.current, "compare", this.url)
      return this.parent.current == this.url
   }

</menu>
