#!/usr/bin/env node
/**
 * Playwright test to launch RyCode TUI and check model selector
 * This will help us debug why models aren't showing
 */

import { spawn } from 'child_process';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

async function testTUI() {
  console.log('Starting TUI test...\n');

  // Clear debug log
  const debugLog = '/tmp/rycode-debug.log';
  if (fs.existsSync(debugLog)) {
    fs.unlinkSync(debugLog);
    console.log('Cleared debug log\n');
  }

  // Launch the TUI
  const tuiPath = path.join(__dirname, 'bin', 'rycode');
  console.log(`Launching TUI: ${tuiPath}\n`);

  const tui = spawn(tuiPath, [], {
    cwd: __dirname,
    stdio: ['pipe', 'pipe', 'pipe'],
    env: { ...process.env, TERM: 'xterm-256color' }
  });

  let output = '';

  tui.stdout.on('data', (data) => {
    output += data.toString();
  });

  tui.stderr.on('data', (data) => {
    console.error('TUI stderr:', data.toString());
  });

  // Wait for TUI to start
  await new Promise(resolve => setTimeout(resolve, 2000));

  console.log('Typing /models...\n');

  // Type /models
  tui.stdin.write('/models');

  // Wait for completion dialog
  await new Promise(resolve => setTimeout(resolve, 500));

  console.log('Pressing Enter to select command...\n');

  // Press Enter to select the "models" command from completion dialog
  tui.stdin.write('\r');

  // Wait for modal to render
  await new Promise(resolve => setTimeout(resolve, 2000));

  console.log('Model dialog should now be open!\n');

  console.log('Closing TUI...\n');

  // Send Escape to close modal
  tui.stdin.write('\x1b');

  await new Promise(resolve => setTimeout(resolve, 500));

  // Send Ctrl+C to quit
  tui.kill('SIGINT');

  // Wait for process to exit
  await new Promise(resolve => {
    tui.on('close', resolve);
    setTimeout(resolve, 1000);
  });

  console.log('\n=== TUI Output (last 500 chars) ===');
  console.log(output.slice(-500));
  console.log('\n=== End TUI Output ===\n');

  // Read debug log
  if (fs.existsSync(debugLog)) {
    console.log('\n=== Debug Log ===');
    const log = fs.readFileSync(debugLog, 'utf8');
    console.log(log);
    console.log('=== End Debug Log ===\n');
  } else {
    console.log('\n⚠️  WARNING: Debug log was not created!\n');
    console.log('This means the TUI binary might not be the latest one or logging failed.\n');
  }

  console.log('\nTest complete!');
}

testTUI().catch(err => {
  console.error('Test failed:', err);
  process.exit(1);
});
