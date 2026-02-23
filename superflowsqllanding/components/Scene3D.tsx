'use client';

import { useEffect, useRef } from 'react';

export default function Scene3D() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Set canvas size
    const resizeCanvas = () => {
      canvas.width = canvas.offsetWidth * window.devicePixelRatio;
      canvas.height = canvas.offsetHeight * window.devicePixelRatio;
      ctx.scale(window.devicePixelRatio, window.devicePixelRatio);
    };
    resizeCanvas();
    window.addEventListener('resize', resizeCanvas);

    const GRID_WIDTH = 40;
    const GRID_DEPTH = 25;
    const GRID_SPACING = 20;
    const WAVE_SPEED = 0.04;
    const WAVE_AMPLITUDE = 60;
    const BALL_COUNT = 8;
    const GRAVITY = 0.5;
    const BOUNCE_DAMPING = 0.7;

    // Grid node class
    class GridNode {
      gridX: number;
      gridZ: number;
      y: number;

      constructor(gridX: number, gridZ: number) {
        this.gridX = gridX;
        this.gridZ = gridZ;
        this.y = 0;
      }

      update(time: number) {
        // Create smooth wave deformation similar to the image
        const centerX = GRID_WIDTH / 2;
        const centerZ = GRID_DEPTH / 2;
        
        // Distance from center
        const distX = (this.gridX - centerX) / centerX;
        const distZ = (this.gridZ - centerZ) / centerZ;
        const dist = Math.sqrt(distX * distX + distZ * distZ);
        
        // Wave pattern
        const wave = Math.sin(dist * 3 - time * WAVE_SPEED) * Math.exp(-dist * 0.5);
        const wave2 = Math.cos(this.gridX * 0.2 + time * WAVE_SPEED * 0.7) * 
                      Math.sin(this.gridZ * 0.15 + time * WAVE_SPEED * 0.9);
        
        this.y = wave * WAVE_AMPLITUDE + wave2 * WAVE_AMPLITUDE * 0.3;
      }

      project(): { x: number; y: number } {
        // Perspective projection - centered and filling the container
        const perspective = 0.7;
        const scale = Math.min(canvas.offsetWidth, canvas.offsetHeight) / (GRID_WIDTH * GRID_SPACING * 0.8);
        const viewAngle = 0.45;
        
        const x = (this.gridX - GRID_WIDTH / 2) * GRID_SPACING * scale;
        const z = (this.gridZ - GRID_DEPTH / 2) * GRID_SPACING * scale;
        
        const depth = z * perspective;
        const perspectiveScale = 1 - Math.abs(z) / (GRID_DEPTH * GRID_SPACING * scale) * 0.2;
        
        const projX = x * perspectiveScale;
        const projY = this.y * scale - depth * viewAngle;
        
        return {
          x: projX + canvas.offsetWidth / 2,
          y: projY + canvas.offsetHeight / 2
        };
      }
    }

    // Ball class
    class Ball {
      gridX: number;
      gridZ: number;
      y: number;
      vy: number;
      radius: number;
      color: string;

      constructor() {
        this.gridX = 5 + Math.random() * (GRID_WIDTH - 10);
        this.gridZ = 5 + Math.random() * (GRID_DEPTH - 10);
        this.y = -150 - Math.random() * 100;
        this.vy = 0;
        this.radius = 10;
        
        const colors = ['#CF4427', '#25A8D6', '#003C6E', '#25A666', '#218ED1'];
        this.color = colors[Math.floor(Math.random() * colors.length)];
      }

      update(time: number) {
        // Apply gravity
        this.vy += GRAVITY;
        this.y += this.vy;

        // Get surface height at ball position
        const surfaceY = this.getSurfaceHeight(time);

        // Check collision with surface
        if (this.y >= surfaceY) {
          this.y = surfaceY;
          this.vy = -Math.abs(this.vy) * BOUNCE_DAMPING;
          
          // Reset if bounce is too weak
          if (Math.abs(this.vy) < 1.5) {
            this.y = -150 - Math.random() * 100;
            this.vy = 0;
            this.gridX = 5 + Math.random() * (GRID_WIDTH - 10);
            this.gridZ = 5 + Math.random() * (GRID_DEPTH - 10);
          }
        }
      }

      getSurfaceHeight(time: number): number {
        const centerX = GRID_WIDTH / 2;
        const centerZ = GRID_DEPTH / 2;
        
        const distX = (this.gridX - centerX) / centerX;
        const distZ = (this.gridZ - centerZ) / centerZ;
        const dist = Math.sqrt(distX * distX + distZ * distZ);
        
        const wave = Math.sin(dist * 3 - time * WAVE_SPEED) * Math.exp(-dist * 0.5);
        const wave2 = Math.cos(this.gridX * 0.2 + time * WAVE_SPEED * 0.7) * 
                      Math.sin(this.gridZ * 0.15 + time * WAVE_SPEED * 0.9);
        
        return wave * WAVE_AMPLITUDE + wave2 * WAVE_AMPLITUDE * 0.3;
      }

      project(): { x: number; y: number; size: number } {
        const perspective = 0.7;
        const scale = Math.min(canvas.offsetWidth, canvas.offsetHeight) / (GRID_WIDTH * GRID_SPACING * 0.8);
        const viewAngle = 0.45;
        
        const x = (this.gridX - GRID_WIDTH / 2) * GRID_SPACING * scale;
        const z = (this.gridZ - GRID_DEPTH / 2) * GRID_SPACING * scale;
        
        const depth = z * perspective;
        const perspectiveScale = 1 - Math.abs(z) / (GRID_DEPTH * GRID_SPACING * scale) * 0.2;
        
        const projX = x * perspectiveScale;
        const projY = this.y * scale - depth * viewAngle;
        
        return {
          x: projX + canvas.offsetWidth / 2,
          y: projY + canvas.offsetHeight / 2,
          size: this.radius * perspectiveScale * scale * 0.8
        };
      }

      draw() {
        const proj = this.project();
        const scale = Math.min(canvas.offsetWidth, canvas.offsetHeight) / (GRID_WIDTH * GRID_SPACING * 0.8);

        // Shadow
        const surfaceY = this.getSurfaceHeight(0);
        const shadowZ = (this.gridZ - GRID_DEPTH / 2) * GRID_SPACING * scale;
        const shadowDepth = shadowZ * 0.7;
        const shadowY = surfaceY * scale - shadowDepth * 0.45;
        
        ctx.fillStyle = 'rgba(0, 0, 0, 0.2)';
        ctx.beginPath();
        ctx.ellipse(proj.x, shadowY + canvas.offsetHeight / 2, proj.size * 0.6, proj.size * 0.3, 0, 0, Math.PI * 2);
        ctx.fill();

        // Ball gradient
        const gradient = ctx.createRadialGradient(
          proj.x - proj.size * 0.25,
          proj.y - proj.size * 0.25,
          0,
          proj.x,
          proj.y,
          proj.size
        );
        gradient.addColorStop(0, this.lightenColor(this.color, 60));
        gradient.addColorStop(0.6, this.color);
        gradient.addColorStop(1, this.darkenColor(this.color, 30));

        ctx.fillStyle = gradient;
        ctx.beginPath();
        ctx.arc(proj.x, proj.y, proj.size, 0, Math.PI * 2);
        ctx.fill();

        // Highlight
        ctx.fillStyle = 'rgba(255, 255, 255, 0.6)';
        ctx.beginPath();
        ctx.arc(proj.x - proj.size * 0.3, proj.y - proj.size * 0.3, proj.size * 0.3, 0, Math.PI * 2);
        ctx.fill();

        // Rim light
        ctx.strokeStyle = 'rgba(255, 255, 255, 0.4)';
        ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.arc(proj.x, proj.y, proj.size - 1, 0, Math.PI * 2);
        ctx.stroke();
      }

      lightenColor(color: string, percent: number): string {
        const num = parseInt(color.replace('#', ''), 16);
        const amt = Math.round(2.55 * percent);
        const R = Math.min(255, (num >> 16) + amt);
        const G = Math.min(255, ((num >> 8) & 0x00FF) + amt);
        const B = Math.min(255, (num & 0x0000FF) + amt);
        return `#${(0x1000000 + R * 0x10000 + G * 0x100 + B).toString(16).slice(1)}`;
      }

      darkenColor(color: string, percent: number): string {
        const num = parseInt(color.replace('#', ''), 16);
        const amt = Math.round(2.55 * percent);
        const R = Math.max(0, (num >> 16) - amt);
        const G = Math.max(0, ((num >> 8) & 0x00FF) - amt);
        const B = Math.max(0, (num & 0x0000FF) - amt);
        return `#${(0x1000000 + R * 0x10000 + G * 0x100 + B).toString(16).slice(1)}`;
      }
    }

    // Create grid
    const gridNodes: GridNode[][] = [];
    for (let x = 0; x < GRID_WIDTH; x++) {
      gridNodes[x] = [];
      for (let z = 0; z < GRID_DEPTH; z++) {
        gridNodes[x][z] = new GridNode(x, z);
      }
    }

    // Create balls
    const balls: Ball[] = [];
    for (let i = 0; i < BALL_COUNT; i++) {
      balls.push(new Ball());
    }

    // Draw grid
    function drawGrid() {
      ctx.strokeStyle = '#003C6E';
      ctx.lineWidth = 1.5;
      ctx.globalAlpha = 0.3;
      ctx.lineCap = 'round';
      ctx.lineJoin = 'round';

      // Draw horizontal lines (along X axis)
      for (let z = 0; z < GRID_DEPTH; z++) {
        ctx.beginPath();
        for (let x = 0; x < GRID_WIDTH; x++) {
          const proj = gridNodes[x][z].project();
          if (x === 0) ctx.moveTo(proj.x, proj.y);
          else ctx.lineTo(proj.x, proj.y);
        }
        ctx.stroke();
      }

      // Draw vertical lines (along Z axis)
      for (let x = 0; x < GRID_WIDTH; x++) {
        ctx.beginPath();
        for (let z = 0; z < GRID_DEPTH; z++) {
          const proj = gridNodes[x][z].project();
          if (z === 0) ctx.moveTo(proj.x, proj.y);
          else ctx.lineTo(proj.x, proj.y);
        }
        ctx.stroke();
      }
      
      ctx.globalAlpha = 1.0;
    }

    // Animation loop
    let time = 0;
    let animationId: number;
    
    const animate = () => {
      time++;
      
      // Clear canvas with transparency to show gradient behind
      ctx.clearRect(0, 0, canvas.offsetWidth, canvas.offsetHeight);

      // Update grid
      for (let x = 0; x < GRID_WIDTH; x++) {
        for (let z = 0; z < GRID_DEPTH; z++) {
          gridNodes[x][z].update(time);
        }
      }

      // Update balls
      balls.forEach(ball => ball.update(time));

      // Draw grid
      drawGrid();

      // Sort and draw balls
      balls.sort((a, b) => a.gridZ - b.gridZ);
      balls.forEach(ball => ball.draw());

      animationId = requestAnimationFrame(animate);
    };

    animate();

    return () => {
      window.removeEventListener('resize', resizeCanvas);
      cancelAnimationFrame(animationId);
    };
  }, []);

  return (
    <div className="w-full h-full relative bg-gradient-to-br from-slate-50 via-blue-50/30 to-slate-50">
      <canvas
        ref={canvasRef}
        className="w-full h-full"
        style={{ display: 'block' }}
        onContextMenu={(e) => e.preventDefault()}
      />
    </div>
  );
}