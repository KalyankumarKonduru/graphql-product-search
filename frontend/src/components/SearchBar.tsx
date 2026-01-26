import { useState, useRef, useEffect } from 'react';
import { Search, X } from 'lucide-react';
import { useDebouncedSearch, useSearchSuggestions } from '../hooks/useProducts';

interface SearchBarProps {
  onSearch: (query: string) => void;
}

export function SearchBar({ onSearch }: SearchBarProps) {
  const { searchTerm, setSearchTerm, debouncedTerm } = useDebouncedSearch(300);
  const { suggestions } = useSearchSuggestions(debouncedTerm);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    onSearch(debouncedTerm);
  }, [debouncedTerm, onSearch]);

  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setShowSuggestions(false);
      }
    };
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, []);

  return (
    <div ref={containerRef} className="relative w-full max-w-xl">
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input
          type="text"
          value={searchTerm}
          onChange={(e) => {
            setSearchTerm(e.target.value);
            setShowSuggestions(true);
          }}
          onFocus={() => setShowSuggestions(true)}
          placeholder="Search products..."
          className="w-full pl-10 pr-10 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent"
        />
        {searchTerm && (
          <button
            onClick={() => setSearchTerm('')}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
          >
            <X className="w-5 h-5" />
          </button>
        )}
      </div>
      
      {showSuggestions && suggestions.length > 0 && searchTerm.length >= 2 && (
        <ul className="absolute z-50 w-full mt-2 bg-white border border-gray-200 rounded-xl shadow-lg overflow-hidden">
          {suggestions.map((suggestion: string) => (
            <li
              key={suggestion}
              onClick={() => {
                setSearchTerm(suggestion);
                setShowSuggestions(false);
              }}
              className="px-4 py-3 hover:bg-gray-50 cursor-pointer flex items-center gap-3"
            >
              <Search className="w-4 h-4 text-gray-400" />
              {suggestion}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
