import React, { useRef, useEffect } from 'react';
import Editor from '@monaco-editor/react';
import { useTheme } from '../context/ThemeContext.jsx';

export default function CodeEditor({ value, onChange, language = 'cpp', onSave, onSubmit, onMount }) {
  const { dark } = useTheme();
  const editorRef = useRef(null);

  const handleMount = (editor) => {
    editorRef.current = editor;
    onMount?.(editor);

    // Register Ctrl+S keybinding in Monaco
    editor.addCommand(
      // Monaco.KeyMod.CtrlCmd | Monaco.KeyCode.KeyS
      2048 | 49, // CtrlCmd = 2048, KeyS = 49
      () => { onSave?.(); }
    );

    // Register Ctrl+Enter keybinding in Monaco
    editor.addCommand(
      // Monaco.KeyMod.CtrlCmd | Monaco.KeyCode.Enter
      2048 | 3, // CtrlCmd = 2048, Enter = 3
      () => { onSubmit?.(); }
    );
  };

  // Global keyboard shortcuts (for when editor is not focused)
  useEffect(() => {
    const handler = (e) => {
      // Ctrl+S / Cmd+S
      if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        e.preventDefault();
        onSave?.();
      }
      // Ctrl+Enter / Cmd+Enter
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault();
        onSubmit?.();
      }
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [onSave, onSubmit]);

  return (
    <Editor
      height="100%"
      language={language}
      value={value}
      theme={dark ? 'vs-dark' : 'vs'}
      onChange={onChange}
      onMount={handleMount}
      options={{
        fontSize: 13,
        fontFamily: "'Fira Code', 'SF Mono', 'Cascadia Code', Menlo, Monaco, Consolas, monospace",
        fontLigatures: true,
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        lineNumbers: 'on',
        renderLineHighlight: 'line',
        tabSize: 4,
        wordWrap: 'off',
        automaticLayout: true,
        padding: { top: 12, bottom: 12 },
        scrollbar: { vertical: 'auto', horizontal: 'auto', verticalScrollbarSize: 6, horizontalScrollbarSize: 6 },
        bracketPairColorization: { enabled: true },
        suggest: { showKeywords: true },
        quickSuggestions: true,
      }}
    />
  );
}
