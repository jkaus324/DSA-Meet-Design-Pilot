// Generic Java test runner driven by spec.yaml + tests/cases/partN.yaml.
// Mirrors the behaviour of harness/python/runner.py byte-for-byte.
//
//   javac -cp "<tmpdir>:<snakeyaml.jar>" Runner.java <user files>
//   java  -cp "<tmpdir>:<snakeyaml.jar>" Runner --problem-dir <tmpdir> --part <N>

import java.io.FileInputStream;
import java.io.InputStream;
import java.lang.reflect.Constructor;
import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

import org.yaml.snakeyaml.Yaml;

public class Runner {

    // ---- CLI ----

    public static void main(String[] argv) throws Exception {
        String problemDir = null;
        Integer part = null;
        for (int i = 0; i < argv.length; i++) {
            switch (argv[i]) {
                case "--problem-dir": problemDir = argv[++i]; break;
                case "--part":        part = Integer.parseInt(argv[++i]); break;
                default: throw new IllegalArgumentException("Unknown arg: " + argv[i]);
            }
        }
        if (problemDir == null || part == null) {
            throw new IllegalArgumentException("--problem-dir and --part are required");
        }

        Yaml yaml = new Yaml();
        Map<String, Object> spec;
        try (InputStream in = new FileInputStream(Paths.get(problemDir, "spec.yaml").toFile())) {
            spec = yaml.load(in);
        }

        Path casesPath = Paths.get(problemDir, "tests", "cases", "part" + part + ".yaml");
        List<Map<String, Object>> cases;
        try (InputStream in = new FileInputStream(casesPath.toFile())) {
            Object raw = yaml.load(in);
            cases = (raw == null) ? new ArrayList<>() : (List<Map<String, Object>>) raw;
        }

        Map<String, Map<String, Object>> fnSpecs = normalizeFunctions(spec);
        Map<String, Map<String, Object>> types = asMap(spec.get("types"));
        Map<String, Map<String, Object>> factories = asMap(spec.get("factories"));

        int passed = 0;
        int total = cases.size();
        for (Map<String, Object> caseMap : cases) {
            String name = (String) caseMap.getOrDefault("name", "<unnamed>");
            boolean ok;
            try {
                ok = runTest(caseMap, fnSpecs, types, factories);
            } catch (Throwable t) {
                ok = false;
                t.printStackTrace(System.err);
            }
            System.out.println((ok ? "PASS " : "FAIL ") + name);
            if (ok) passed++;
        }
        System.out.println("PART" + part + "_SUMMARY " + passed + "/" + total);
    }

    // ---- Spec helpers ----

    @SuppressWarnings("unchecked")
    private static Map<String, Map<String, Object>> normalizeFunctions(Map<String, Object> spec) {
        Object fns = spec.get("functions");
        Map<String, Map<String, Object>> out = new LinkedHashMap<>();
        if (fns == null) return out;
        if (fns instanceof List) {
            for (Object e : (List<Object>) fns) {
                Map<String, Object> m = (Map<String, Object>) e;
                out.put((String) m.get("name"), m);
            }
        } else if (fns instanceof Map) {
            for (Map.Entry<String, Object> e : ((Map<String, Object>) fns).entrySet()) {
                out.put(e.getKey(), (Map<String, Object>) e.getValue());
            }
        }
        return out;
    }

    @SuppressWarnings("unchecked")
    private static <V> Map<String, V> asMap(Object o) {
        return (o == null) ? new LinkedHashMap<>() : (Map<String, V>) o;
    }

    private static String parseListType(Object t) {
        if (t instanceof String s && s.startsWith("list<") && s.endsWith(">")) {
            return s.substring(5, s.length() - 1);
        }
        return null;
    }

    // ---- Deserialization ----

    @SuppressWarnings("unchecked")
    private static Object deserializeArg(
            Object value,
            String declaredType,
            Map<String, Map<String, Object>> types,
            Map<String, Map<String, Object>> factories) throws Exception {

        String inner = parseListType(declaredType);
        if (inner != null) {
            if (!(value instanceof List)) {
                throw new IllegalArgumentException("expected list for " + declaredType);
            }
            List<Object> out = new ArrayList<>();
            for (Object item : (List<Object>) value) {
                out.add(deserializeArg(item, inner, types, factories));
            }
            return out;
        }

        if ("factory".equals(declaredType)) {
            Map<String, Object> fac = factories.get((String) value);
            if (fac == null) throw new IllegalArgumentException("unknown factory: " + value);
            String expr = (String) fac.get("java");
            return evalJavaFactory(expr);
        }

        if (declaredType != null && types.containsKey(declaredType) && value instanceof Map) {
            return deserializeStruct((Map<String, Object>) value, declaredType, types);
        }

        // Primitives / strings / unknown — pass through.
        return value;
    }

    @SuppressWarnings("unchecked")
    private static Object deserializeStruct(
            Map<String, Object> value,
            String typeName,
            Map<String, Map<String, Object>> types) throws Exception {

        Map<String, Object> typeDef = types.get(typeName);
        List<Map<String, Object>> fields = (List<Map<String, Object>>) typeDef.get("fields");

        List<Object> args = new ArrayList<>();
        // Build full positional argument list: every field gets either the
        // YAML value or the spec-declared default. This mirrors Python kwarg
        // semantics so callers can omit any optional field freely.
        for (Map<String, Object> f : fields) {
            String fname = (String) f.get("name");
            if (value.containsKey(fname)) {
                args.add(value.get(fname));
            } else if (f.containsKey("default")) {
                args.add(f.get("default"));
            } else {
                throw new IllegalArgumentException(
                        "missing required field '" + fname + "' for " + typeName);
            }
        }

        Class<?> cls = Class.forName(typeName);
        Constructor<?> ctor = findConstructor(cls, args.size());
        Object[] coerced = new Object[args.size()];
        Class<?>[] params = ctor.getParameterTypes();
        for (int i = 0; i < args.size(); i++) {
            coerced[i] = coerce(args.get(i), params[i]);
        }
        return ctor.newInstance(coerced);
    }

    /** Parse a tiny subset: `new ClassName()` or `new ClassName(args...)`. */
    private static Object evalJavaFactory(String expr) throws Exception {
        String s = expr.trim();
        if (!s.startsWith("new ")) {
            throw new IllegalArgumentException("Unsupported factory expr: " + expr);
        }
        s = s.substring(4).trim();
        int open = s.indexOf('(');
        int close = s.lastIndexOf(')');
        if (open < 0 || close < 0 || close < open) {
            throw new IllegalArgumentException("Malformed factory expr: " + expr);
        }
        String className = s.substring(0, open).trim();
        String inside = s.substring(open + 1, close).trim();

        Class<?> cls = Class.forName(className);
        if (inside.isEmpty()) {
            return cls.getDeclaredConstructor().newInstance();
        }
        // Minimal arg parsing — split on top-level commas, parse each as int/double/bool/string literal.
        List<String> rawArgs = splitTopLevelCommas(inside);
        Object[] vals = new Object[rawArgs.size()];
        Class<?>[] sig = new Class<?>[rawArgs.size()];
        for (int i = 0; i < rawArgs.size(); i++) {
            Object v = parseLiteral(rawArgs.get(i).trim());
            vals[i] = v;
            sig[i] = (v == null) ? Object.class : primitiveSig(v.getClass());
        }
        Constructor<?> ctor = findCtorByArity(cls, vals.length, sig);
        Object[] coerced = new Object[vals.length];
        Class<?>[] params = ctor.getParameterTypes();
        for (int i = 0; i < vals.length; i++) coerced[i] = coerce(vals[i], params[i]);
        return ctor.newInstance(coerced);
    }

    private static List<String> splitTopLevelCommas(String s) {
        List<String> out = new ArrayList<>();
        int depth = 0, start = 0;
        boolean inStr = false;
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            if (c == '"' && (i == 0 || s.charAt(i - 1) != '\\')) inStr = !inStr;
            else if (!inStr && (c == '(' || c == '{' || c == '[')) depth++;
            else if (!inStr && (c == ')' || c == '}' || c == ']')) depth--;
            else if (!inStr && c == ',' && depth == 0) {
                out.add(s.substring(start, i));
                start = i + 1;
            }
        }
        out.add(s.substring(start));
        return out;
    }

    private static Object parseLiteral(String s) {
        if (s.equals("true")) return Boolean.TRUE;
        if (s.equals("false")) return Boolean.FALSE;
        if (s.equals("null")) return null;
        if (s.length() >= 2 && s.charAt(0) == '"' && s.charAt(s.length() - 1) == '"') {
            return s.substring(1, s.length() - 1);
        }
        try { return Integer.parseInt(s); } catch (NumberFormatException ignore) {}
        try { return Long.parseLong(s); } catch (NumberFormatException ignore) {}
        try { return Double.parseDouble(s); } catch (NumberFormatException ignore) {}
        return s;
    }

    private static Class<?> primitiveSig(Class<?> c) {
        if (c == Integer.class) return int.class;
        if (c == Long.class) return long.class;
        if (c == Double.class) return double.class;
        if (c == Float.class) return float.class;
        if (c == Boolean.class) return boolean.class;
        return c;
    }

    private static Constructor<?> findConstructor(Class<?> cls, int arity) {
        Constructor<?> best = null;
        for (Constructor<?> c : cls.getDeclaredConstructors()) {
            if (c.getParameterCount() == arity) {
                c.setAccessible(true);
                if (best == null || c.getParameterCount() < best.getParameterCount()) best = c;
            }
        }
        if (best == null) {
            throw new IllegalStateException(
                    "No constructor with arity " + arity + " for " + cls.getName());
        }
        return best;
    }

    private static Constructor<?> findCtorByArity(Class<?> cls, int arity, Class<?>[] sigHint) {
        // Prefer exact arity; types are coerced afterward.
        return findConstructor(cls, arity);
    }

    /** Coerce YAML-loaded value into the constructor / setter parameter type. */
    private static Object coerce(Object v, Class<?> target) {
        if (v == null) return null;
        if (target.isInstance(v)) return v;
        if (target == int.class || target == Integer.class) {
            if (v instanceof Number n) return n.intValue();
        }
        if (target == long.class || target == Long.class) {
            if (v instanceof Number n) return n.longValue();
        }
        if (target == double.class || target == Double.class) {
            if (v instanceof Number n) return n.doubleValue();
        }
        if (target == float.class || target == Float.class) {
            if (v instanceof Number n) return n.floatValue();
        }
        if (target == boolean.class || target == Boolean.class) {
            if (v instanceof Boolean b) return b;
        }
        if (target == String.class && v instanceof String s) return s;
        return v;
    }

    // ---- Test execution ----

    @SuppressWarnings("unchecked")
    private static boolean runTest(
            Map<String, Object> caseMap,
            Map<String, Map<String, Object>> fnSpecs,
            Map<String, Map<String, Object>> types,
            Map<String, Map<String, Object>> factories) throws Exception {

        String fnName = (String) caseMap.get("call");
        Map<String, Object> fnSpec = fnSpecs.get(fnName);
        List<Map<String, Object>> params = (List<Map<String, Object>>) fnSpec.get("params");

        List<Object> rawArgs = (List<Object>) caseMap.getOrDefault("args", new ArrayList<>());
        Object[] args = new Object[rawArgs.size()];
        for (int i = 0; i < rawArgs.size(); i++) {
            args[i] = deserializeArg(rawArgs.get(i), (String) params.get(i).get("type"), types, factories);
        }

        Class<?> solCls = Class.forName("Solution");
        Method method = findStaticMethod(solCls, fnName, args.length);
        // Coerce args against the method signature.
        Class<?>[] sig = method.getParameterTypes();
        Object[] coerced = new Object[args.length];
        for (int i = 0; i < args.length; i++) coerced[i] = coerce(args[i], sig[i]);

        boolean expectThrows = Boolean.TRUE.equals(caseMap.get("expect_throws"));
        Object result;
        try {
            result = method.invoke(null, coerced);
        } catch (InvocationTargetException ite) {
            return expectThrows;
        } catch (Throwable t) {
            return expectThrows;
        }
        if (expectThrows) return false;

        if (caseMap.containsKey("expect_equals")) {
            Object expected = caseMap.get("expect_equals");
            if (!scalarEquals(result, expected)) return false;
        }

        if (caseMap.containsKey("expect_close")) {
            double eps = 0.001;
            if (caseMap.containsKey("epsilon")) {
                eps = ((Number) caseMap.get("epsilon")).doubleValue();
            }
            double expected = ((Number) caseMap.get("expect_close")).doubleValue();
            if (!(result instanceof Number)) return false;
            double actual = ((Number) result).doubleValue();
            if (Math.abs(actual - expected) > eps) return false;
        }

        if (caseMap.containsKey("expect_size")) {
            int expectedSize = ((Number) caseMap.get("expect_size")).intValue();
            int actualSize = sizeOf(result);
            if (actualSize != expectedSize) return false;
        }

        if (caseMap.containsKey("expect")) {
            String field = (String) caseMap.get("expect_field");
            List<Object> expected = (List<Object>) caseMap.get("expect");
            List<Object> actual = new ArrayList<>();
            Iterator<?> it = iter(result);
            while (it.hasNext()) {
                Object item = it.next();
                actual.add(field == null ? item : extractField(item, field));
            }
            if (!listEquals(actual, expected)) return false;
        }

        Object alsoRaw = caseMap.get("also");
        if (alsoRaw instanceof List) {
            List<Object> resultList = toList(result);
            for (Object check : (List<Object>) alsoRaw) {
                Map<String, Object> c = (Map<String, Object>) check;
                int idx = ((Number) c.get("index")).intValue();
                String field = (String) c.get("field");
                Object expected = c.get("equals");
                Object actual = extractField(resultList.get(idx), field);
                if (!scalarEquals(actual, expected)) return false;
            }
        }

        return true;
    }

    private static Method findStaticMethod(Class<?> cls, String name, int arity) {
        for (Method m : cls.getDeclaredMethods()) {
            if (m.getName().equals(name) && m.getParameterCount() == arity) {
                m.setAccessible(true);
                return m;
            }
        }
        throw new IllegalStateException(
                "No method " + name + "/" + arity + " on " + cls.getName());
    }

    private static int sizeOf(Object o) {
        if (o instanceof List<?> l) return l.size();
        if (o instanceof Object[] a) return a.length;
        if (o instanceof Iterable<?> it) {
            int n = 0;
            for (Object ignore : it) n++;
            return n;
        }
        throw new IllegalArgumentException("Not sizeable: " + (o == null ? "null" : o.getClass()));
    }

    @SuppressWarnings("unchecked")
    private static List<Object> toList(Object o) {
        if (o instanceof List<?> l) return (List<Object>) l;
        if (o instanceof Object[] a) return new ArrayList<>(Arrays.asList(a));
        if (o instanceof Iterable<?> it) {
            List<Object> out = new ArrayList<>();
            for (Object item : it) out.add(item);
            return out;
        }
        throw new IllegalArgumentException("Not iterable: " + (o == null ? "null" : o.getClass()));
    }

    private static Iterator<?> iter(Object o) {
        return toList(o).iterator();
    }

    /** Try `obj.field` then `obj.getField()`. */
    private static Object extractField(Object obj, String field) throws Exception {
        if (obj instanceof Map<?, ?> m) return m.get(field);
        Class<?> cls = obj.getClass();
        try {
            Field f = cls.getField(field);
            return f.get(obj);
        } catch (NoSuchFieldException ignore) {}
        try {
            Field f = cls.getDeclaredField(field);
            f.setAccessible(true);
            return f.get(obj);
        } catch (NoSuchFieldException ignore) {}
        String getter = "get" + Character.toUpperCase(field.charAt(0)) + field.substring(1);
        try {
            Method m = cls.getMethod(getter);
            return m.invoke(obj);
        } catch (NoSuchMethodException ignore) {}
        // boolean isFoo()
        String isGetter = "is" + Character.toUpperCase(field.charAt(0)) + field.substring(1);
        try {
            Method m = cls.getMethod(isGetter);
            return m.invoke(obj);
        } catch (NoSuchMethodException ignore) {}
        throw new NoSuchFieldException(field + " on " + cls.getName());
    }

    private static boolean listEquals(List<Object> a, List<Object> b) {
        if (a.size() != b.size()) return false;
        for (int i = 0; i < a.size(); i++) {
            if (!scalarEquals(a.get(i), b.get(i))) return false;
        }
        return true;
    }

    private static boolean scalarEquals(Object a, Object b) {
        if (a == b) return true;
        if (a == null || b == null) return false;
        if (a instanceof Number na && b instanceof Number nb) {
            // Compare as double when at least one side is float-y, else as long.
            if (a instanceof Double || a instanceof Float || b instanceof Double || b instanceof Float) {
                return na.doubleValue() == nb.doubleValue();
            }
            return na.longValue() == nb.longValue();
        }
        if (a instanceof List<?> la && b instanceof List<?> lb) {
            if (la.size() != lb.size()) return false;
            for (int i = 0; i < la.size(); i++) {
                if (!scalarEquals(la.get(i), lb.get(i))) return false;
            }
            return true;
        }
        return Objects.equals(a, b);
    }
}
