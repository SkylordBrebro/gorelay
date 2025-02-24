## Package Structure

### Server Components
- `pkg/server/local.go` - TCP server implementation for handling game connections
  - Manages client connections and packet handling
  - Implements connection pooling and concurrent client handling
  - Provides low-level network communication layer
- `pkg/server/monitor.go` - HTTP server for monitoring and debugging purposes
  - Exposes monitoring endpoints (/status, /clients)
  - Provides real-time server statistics and diagnostics
  - Supports custom handler registration for extensibility

### Core Packages
- `pkg/account` - Account management and authentication functionality
- `pkg/client` - Client implementation for game connections
- `pkg/config` - Configuration management and settings
- `pkg/crypto` - Cryptographic utilities including RC4 and RSA implementations
  - RC4 stream cipher for packet encryption/decryption
  - Separate inbound/outbound cipher streams
  - Secure key management
- `pkg/events` - Event system for handling game events
  - Event types for player, enemy, and game actions
  - Subscription-based event handling
  - Type-safe event data structures
  - Support for custom event handlers
- `pkg/logger` - Logging system for application-wide logging
- `pkg/models` - Core data models including:
  - Game entities and objects:
    - Base Entity type with position, size, and condition tracking
    - Player entities with inventory, stats, and class information
    - Enemy entities with HP, defense, and behavior flags
    - Containers for inventory management
    - Projectiles with damage and trajectory data
    - Portals with dungeon and state information
  - Character classes
  - Map information
  - Account structures
  - Guild systems
  - Pet systems
  - Movement records
  - Game IDs and endpoints
- `pkg/packets` - Network packet definitions and handling
  - Core packet types for game state synchronization:
    - Movement and position updates
    - Combat actions (shooting, AOE effects)
    - World state updates (new ticks, map info)
    - Chat and text messages
    - Inventory and item usage
  - Binary encoding/decoding for network transmission
  - Packet ID system for message routing
- `pkg/plugin` - Plugin system for extending functionality
  - Dynamic plugin loading and lifecycle management
  - Packet hook registration system
  - Plugin interface with Initialize/Enable/Disable methods
  - Support for runtime plugin loading and unloading
  - Automatic packet handler registration
- `pkg/resources` - Resource management and assets
  - Game object definitions:
    - Objects with properties (HP, defense, equipment slots)
    - Tiles with behavior flags (walkable, damage, speed)
    - Pets with abilities and families
    - Projectiles with movement patterns and damage
  - JSON-based resource loading
  - In-memory caching of game definitions
  - Type-safe access to game resources
- `pkg/services` - Utility services including:
  - HTTP client functionality
  - Pathfinding algorithms:
    - A* implementation for optimal path calculation
    - Dynamic walkability updates
    - Efficient node management with priority queue
    - Support for diagonal movement
  - Random number generation
  - String utilities
  - XML to JSON conversion
  - Update management

### Implementation Details
The server architecture uses a dual-server approach:
1. A TCP server (`LocalServer`) that handles real-time game connections with:
   - Concurrent client handling
   - Binary packet processing
   - Connection state management
2. A monitoring HTTP server (`MonitorServer`) that provides:
   - Real-time server statistics
   - Client connection tracking
   - Diagnostic endpoints
   - Extensible monitoring capabilities

The game world is built around a robust entity system where all game objects inherit from a base `Entity` type, providing:
- Consistent position and state tracking
- Condition effect management
- Movement timestamping

### Network Protocol
The game uses a binary packet-based protocol with:
1. Core packet types:
   - World state updates (NewTickPacket, UpdatePacket)
   - Player actions (MovePacket, PlayerShootPacket)
   - Enemy behavior (EnemyShootPacket, AoePacket)
   - Item management (UseItemPacket)
   - Communication (TextPacket)
2. Packet structure:
   - Unique packet IDs for routing
   - Position tracking using WorldPosData
   - Stat management with StatData
   - Object state synchronization
3. State management:
   - Tick-based updates
   - Real-time movement recording
   - Server-client time synchronization

### Plugin System
The plugin architecture provides:
1. Dynamic loading:
   - Runtime plugin loading and unloading
   - Plugin lifecycle management (Initialize/Enable/Disable)
   - Automatic cleanup on unload
2. Packet hooks:
   - Automatic method discovery and registration
   - Type-safe packet handling
   - Multiple hooks per packet type
3. Plugin interface:
   - Standard metadata (Name, Author, Version)
   - Client instance access
   - Error handling and initialization

### Resource Management
The resource system handles:
1. Game definitions:
   - Objects with complete property sets
   - Tile types with behavior flags
   - Pet definitions with abilities
   - Projectile patterns and behaviors
2. Loading and caching:
   - JSON-based resource files
   - In-memory caching for performance
   - Type-safe access methods
3. Properties:
   - Object stats (HP, defense)
   - Tile behaviors (walkable, damage)
   - Pet abilities and families
   - Projectile patterns (wavy, parametric, boomerang)

### Event System
The event system provides:
1. Event types:
   - Player events (join, leave, move, shoot, hit, death)
   - Enemy events (spawn, death, shoot)
   - Game events (map change, tick, chat)
2. Event handling:
   - Subscribe/unsubscribe mechanism
   - Type-safe event data
   - Asynchronous event dispatch
3. Event data structures:
   - Player event data (position, stats)
   - Enemy event data (type, position)
   - Chat event data (name, message, recipient)

### Pathfinding
The A* pathfinding implementation includes:
1. Core features:
   - Optimal path calculation
   - Support for diagonal movement
   - Dynamic walkability updates
2. Performance optimizations:
   - Priority queue for efficient node selection
   - Node reuse to minimize allocations
   - Cached calculations
3. Customization:
   - Configurable heuristics
   - Walkability updates
   - Grid size management

### Application Configuration
The application uses a structured configuration approach:
1. Command-line flags:
   - Config file path
   - Accounts file path
   - Debug mode toggle
2. Configuration files:
   - JSON-based configuration
   - Account credentials and settings
   - Plugin configuration
3. Resource files:
   - Game object definitions
   - Tile properties
   - Resource caching

### Startup Process
The application initialization follows these steps:
1. Configuration loading:
   - Parse command-line flags
   - Load configuration files
   - Initialize logging system
2. Resource initialization:
   - Load account information
   - Load game resources
   - Initialize resource managers
3. Client setup:
   - Create client instances
   - Load and configure plugins
   - Establish connections
4. Shutdown handling:
   - Graceful shutdown on signals
   - Client disconnection
   - Resource cleanup