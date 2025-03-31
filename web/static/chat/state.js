class State {
    #currentUser = null;  // Private field
    #users = new Map();   // Private field
    #subscribers = [];    // Private field
    #lastAuthCheck = null; // Private field to track when auth was last checked
    #authCacheExpiryMs = 30000; // Auth cache expires after 30 seconds by default
  
    getUsers() {
      return Array.from(this.#users.values());
    }
  
    getUser(id) {
      return this.#users.get(id) || null;
    }
  
    addUser(user) {
      if (!user || typeof user !== 'object') {
        throw new Error('Invalid user object');
      }
      if (user.id === undefined || user.id === null) {
        throw new Error('User must have a valid id property');
      }
      if (this.#users.has(user.id)) {
        throw new Error(`User with id ${user.id} already exists`);
      }
      this.#users.set(user.id, user);
    }
  
    updateUser(id, updatedProperties) {
      if (!this.#users.has(id)) {
        throw new Error(`User with id ${id} does not exist`);
      }
      if (!updatedProperties || typeof updatedProperties !== 'object') {
        throw new Error('updatedProperties must be an object');
      }
      const user = this.#users.get(id);
      const updatedUser = { ...user, ...updatedProperties };
      this.#users.set(id, updatedUser);
    }
  
    removeUser(id) {
      if (!this.#users.has(id)) {
        throw new Error(`User with id ${id} does not exist`);
      }
      const user = this.#users.get(id);
      this.#users.delete(id);
      this._notifySubscribers({ type: 'REMOVE_USER', payload: user });
  
      // If the removed user was the authenticated user, log them out.
      if (this.#currentUser && this.#currentUser.id === id) {
        this.logout();
      }
    }
  
    getCurrentUser() {
      return this.#currentUser;
    }

    setCurrentUser(user) {
        console.log("CURRENT USER: ", user)
        this.#currentUser = user;
    }
  
    isAuthenticated() {
      return this.#currentUser !== null;
    }
  }
  